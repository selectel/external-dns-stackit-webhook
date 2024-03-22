package selprovider

import (
	"context"
	"fmt"

	domains "github.com/selectel/domains-go/pkg/v2"
	"go.uber.org/zap"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
)

// ApplyChanges applies a given set of changes.
func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	client, err := p.getDomainsClient()
	if err != nil {
		return err
	}

	// create rr set. POST /zones/{zoneId}/rrset
	err = p.createRRSets(ctx, client, changes.Create)
	if err != nil {
		return err
	}

	// update rr set. PATCH /zones/{zoneId}/rrset/{rrSetId}
	err = p.updateRRSets(ctx, client, changes.UpdateNew)
	if err != nil {
		return err
	}

	// delete rr set. DELETE /zones/{zoneId}/rrset/{rrSetId}
	err = p.deleteRRSets(ctx, client, changes.Delete)
	if err != nil {
		return err
	}

	return nil
}

// createRRSets creates new record sets for the given endpoints that are in the creation field.
func (p *Provider) createRRSets(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	endpoints []*endpoint.Endpoint,
) error {
	if len(endpoints) == 0 {
		return nil
	}

	return p.handleRRSetWithWorkers(ctx, client, endpoints, CREATE)
}

// updateRRSets patches (overrides) contents in the record sets for the given endpoints that are in the update new field.
func (p *Provider) updateRRSets(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	endpoints []*endpoint.Endpoint,
) error {
	if len(endpoints) == 0 {
		return nil
	}

	return p.handleRRSetWithWorkers(ctx, client, endpoints, UPDATE)
}

// deleteRRSets delete record sets for the given endpoints that are in the deletion field.
func (p *Provider) deleteRRSets(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	endpoints []*endpoint.Endpoint,
) error {
	if len(endpoints) == 0 {
		p.logger.Debug("no endpoints to delete")

		return nil
	}

	p.logger.Info("records to delete", zap.String("records", fmt.Sprintf("%v", endpoints)))

	return p.handleRRSetWithWorkers(ctx, client, endpoints, DELETE)
}

// handleRRSetWithWorkers handles the given endpoints with workers to optimize speed.
func (p *Provider) handleRRSetWithWorkers(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	endpoints []*endpoint.Endpoint,
	action string,
) error {
	zones, err := p.zoneFetcherClient.zones(ctx, client)
	if err != nil {
		return err
	}

	workerChannel := make(chan changeTask, len(endpoints))
	defer close(workerChannel)
	errorChannel := make(chan error, len(endpoints))

	for i := 0; i < p.workers; i++ {
		go p.changeWorker(ctx, client, workerChannel, errorChannel, zones)
	}

	for _, change := range endpoints {
		workerChannel <- changeTask{
			action: action,
			change: change,
		}
	}

	for i := 0; i < len(endpoints); i++ {
		err := <-errorChannel
		if err != nil {
			return err
		}
	}

	return nil
}

// createRRSet creates a new record set for the given endpoint.
func (p *Provider) createRRSet(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	change *endpoint.Endpoint,
	zones []*domains.Zone,
) error {
	resultZone, found := findBestMatchingZone(change.DNSName, zones)
	if !found {
		return fmt.Errorf("no matching zone found for %s", change.DNSName)
	}

	logFields := getLogFields(change, CREATE, resultZone.ID)
	p.logger.Info("create record set", logFields...)

	if p.dryRun {
		p.logger.Debug("dry run, skipping", logFields...)

		return nil
	}

	modifyChange(change)

	rrSet := getRRSetRecord(change)

	// ignore all errors to just retry on next run
	_, err := client.CreateRRSet(ctx, resultZone.ID, rrSet)
	if err != nil {
		p.logger.Error("error creating record set", zap.Error(err))

		return err
	}

	p.logger.Info("create record set successfully", logFields...)

	return nil
}

// updateRRSet patches (overrides) contents in the record set.
func (p *Provider) updateRRSet(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	change *endpoint.Endpoint,
	zones []*domains.Zone,
) error {
	modifyChange(change)

	resultZone, resultRRSet, err := p.rrSetFetcherClient.getRRSetForUpdateDeletion(ctx, client, change, zones)
	if err != nil {
		return err
	}

	logFields := getLogFields(change, UPDATE, resultRRSet.ID)
	p.logger.Info("update record set", logFields...)

	if p.dryRun {
		p.logger.Debug("dry run, skipping", logFields...)

		return nil
	}

	rrSet := getRRSetRecord(change)

	err = client.UpdateRRSet(ctx, resultZone.ID, resultRRSet.ID, rrSet)
	if err != nil {
		p.logger.Error("error updating record set", zap.Error(err))

		return err
	}

	p.logger.Info("record set updated successfully", logFields...)

	return nil
}

// deleteRRSet deletes a record set for the given endpoint.
func (p *Provider) deleteRRSet(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	change *endpoint.Endpoint,
	zones []*domains.Zone,
) error {
	modifyChange(change)

	resultZone, resultRRSet, err := p.rrSetFetcherClient.getRRSetForUpdateDeletion(ctx, client, change, zones)
	if err != nil {
		return err
	}

	logFields := getLogFields(change, DELETE, resultRRSet.ID)
	p.logger.Info("delete record set", logFields...)

	if p.dryRun {
		p.logger.Debug("dry run, skipping", logFields...)

		return nil
	}

	err = client.DeleteRRSet(ctx, resultZone.ID, resultRRSet.ID)
	if err != nil {
		p.logger.Error("error deleting record set", zap.Error(err))

		return err
	}

	p.logger.Info("delete record set successfully", logFields...)

	return nil
}

// changeWorker is a worker that handles changes passed by a channel.
func (p *Provider) changeWorker(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	changes chan changeTask,
	errorChannel chan error,
	zones []*domains.Zone,
) {
	for change := range changes {
		switch change.action {
		case CREATE:
			err := p.createRRSet(ctx, client, change.change, zones)
			errorChannel <- err
		case UPDATE:
			err := p.updateRRSet(ctx, client, change.change, zones)
			errorChannel <- err
		case DELETE:
			err := p.deleteRRSet(ctx, client, change.change, zones)
			errorChannel <- err
		}
	}

	p.logger.Debug("change worker finished")
}
