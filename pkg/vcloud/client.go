package vcloud

import (
	"errors"

	"github.com/proact-de/vcloud-csi-driver/pkg/model"
	"github.com/rs/zerolog/log"
	"github.com/vmware/go-vcloud-director/types/v56"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

// Client is a wrapper around the upstream vCloud Director client.
type Client struct {
	upstream *govcd.VCDClient
	org      *govcd.Org
	vdc      *govcd.Vdc

	username     string
	password     string
	organization string
	datacenter   string
}

// NewClient simply initializes a new client wrapper.
func NewClient(opts ...Option) (*Client, error) {
	options := newOptions(opts...)

	client := &Client{
		upstream: govcd.NewVCDClient(
			*options.Href,
			options.Insecure,
		),
		username:     options.Username,
		password:     options.Password,
		organization: options.Organization,
		datacenter:   options.Datacenter,
	}

	if err := client.Authenticate(); err != nil {
		return nil, err
	}

	{
		org, err := client.upstream.GetOrgByNameOrId(options.Organization)

		if err != nil {
			return nil, err
		}

		client.org = org
	}

	{
		vdc, err := client.org.GetVDCByName(options.Datacenter, true)

		if err != nil {
			return nil, err
		}

		client.vdc = vdc
	}

	return client, nil
}

// Datacenter just returns the name of the vDC.
func (c *Client) Datacenter() string {
	return c.vdc.Vdc.Name
}

// Authenticate wraps the auth for the cloud provider.
func (c *Client) Authenticate() error {
	return c.upstream.Authenticate(
		c.username,
		c.password,
		c.organization,
	)
}

// Disconnect wraps the logout for the cloud provider.
func (c *Client) Disconnect() error {
	return c.upstream.Disconnect()
}

// Refresh simply refreshs the authentication to the API.
func (c *Client) Refresh() error {
	resp, err := c.upstream.GetAuthResponse(
		c.username,
		c.password,
		c.organization,
	)

	if err != nil {
		return err
	}

	if resp.StatusCode == 401 {
		return c.Authenticate()
	}

	return nil
}

// List fetches all available disks.
func (c *Client) List() ([]*model.Volume, error) {
	log.Info().
		Msg("Listing volumes")

	if err := c.Refresh(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to refresh authentication")

		return nil, err
	}

	if err := c.vdc.Refresh(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to refresh datacenter")

		return nil, err
	}

	resp := make([]*model.Volume, 0)

	for _, entity := range c.vdc.Vdc.ResourceEntities {
		for _, resource := range entity.ResourceEntity {
			if resource.Type == types.MimeDisk {
				disk, err := c.vdc.GetDiskByHref(resource.HREF)

				if err != nil {
					return nil, err
				}

				log.Info().
					Interface("disk", disk).
					Msg("")

				resp = append(resp, &model.Volume{
					ID:      resource.ID,
					Name:    resource.Name,
					Size:    disk.Disk.Size,
					Profile: disk.Disk.StorageProfile.Name,
					BusType: disk.Disk.BusType,
					SubType: disk.Disk.BusSubType,
				})
			}
		}
	}

	return resp, nil
}

// Find searches a disk by it's ID.
func (c *Client) Find(name string) (*model.Volume, error) {
	log.Info().
		Str("name", name).
		Msg("Finding volume")

	if err := c.Refresh(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to refresh authentication")

		return nil, err
	}

	if err := c.vdc.Refresh(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to refresh datacenter")

		return nil, err
	}

	disk, err := c.vdc.GetDiskById(name, false)

	if err != nil {
		return nil, err
	}

	return &model.Volume{
		Name:    disk.Disk.Name,
		Size:    disk.Disk.Size,
		Profile: disk.Disk.StorageProfile.Name,
		BusType: disk.Disk.BusType,
		SubType: disk.Disk.BusSubType,
	}, nil
}

// Create is used to create an independent disk.
func (c *Client) Create(name string, capacity int64, profile string) error {
	log.Info().
		Str("name", name).
		Int64("capacity", capacity).
		Str("profile", profile).
		Msg("Creating volume")

	if err := c.Refresh(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to refresh authentication")

		return err
	}

	if err := c.vdc.Refresh(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to refresh datacenter")

		return err
	}

	exists, err := c.exists(name)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return nil
}

func (c *Client) exists(name string) (bool, error) {
	disks, err := c.vdc.GetDisksByName(name, false)

	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			return false, nil
		}

		return false, err
	}

	for _, disk := range *disks {
		log.Debug().Interface("disk", disk).Msg("Found disk")
		return true, nil
	}

	return false, nil
}
