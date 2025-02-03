package bcc_terraform

import (
	"fmt"
	"strings"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token            string
	APIEndpoint      string
	TerraformVersion string
	ClientID         string
}

type CombinedConfig struct {
	manager *bcc.Manager
}

func (c *CombinedConfig) Manager() *bcc.Manager { return c.manager }

func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	manager := bcc.NewManager(c.Token)
	manager.Logger = logger
	manager.BaseURL = strings.TrimSuffix(c.APIEndpoint, "/")
	manager.ClientID = c.ClientID
	manager.UserAgent = fmt.Sprintf("Terraform/%s", c.TerraformVersion)

	return &CombinedConfig{
		manager: manager,
	}, nil
}
