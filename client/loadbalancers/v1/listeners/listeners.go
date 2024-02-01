package listeners

import (
	"fmt"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var (
	listenerIDText = "listener_id is mandatory argument"
	protocolTypes  = types.ProtocolType("").StringList()
)

var listenerListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "loadbalancer listeners list",
	Category: "listener",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "loadbalancer-id",
			Aliases:     []string{"l"},
			Usage:       "loadbalancer ID",
			Required:    false,
			DefaultText: "<nil>",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewLBListenerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := listeners.ListOpts{LoadBalancerID: utils.StringToPointer(c.String("loadbalancer-id"))}

		results, err := listeners.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var listenerCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "create loadbalancer listener",
	Category: "listener",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "listener name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "port",
			Aliases:  []string{"p"},
			Usage:    "listener port",
			Value:    80,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "loadbalancer-id",
			Aliases:  []string{"l"},
			Usage:    "loadbalancer ID",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "secret-id",
			Aliases: []string{"s"},
			Usage:   "ID of the secret where PKCS12 file is stored for TERMINATED_HTTPS load balancer",
		},
		&cli.StringSliceFlag{
			Name:  "sni-secret-id",
			Usage: "List of secret's ID containing PKCS12 format certificate/key bundles for TERMINATED_HTTPS listeners",
		},
		&cli.GenericFlag{
			Name:    "protocol-type",
			Aliases: []string{"pt"},
			Value: &utils.EnumValue{
				Enum: protocolTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypes, ", ")),
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "allowed-cidrs",
			Usage: "List of networks from which listener is accessible",
		},
		&cli.IntFlag{
			Name:     "timeout-client-data",
			Aliases:  []string{"tcd"},
			Usage:    "Frontend client inactivity timeout in milliseconds",
			Value: 50000,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-member-connect",
			Aliases:  []string{"tmc"},
			Usage:    "Backend member connection timeout in milliseconds",
			Value: 5000,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-member-data",
			Aliases:  []string{"tmd"},
			Usage:    "Backend member inactivity timeout in milliseconds",
			Value: 50000,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "connection-limit",
			Aliases:  []string{"cl"},
			Usage:    "Limit of the simultaneous connections",
			Value: 100000,
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewLBListenerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		pt := types.ProtocolType(c.String("protocol-type"))
		if err := pt.IsValid(); err != nil {
			return cli.NewExitError(err, 1)
		}

		opts := listeners.CreateOpts{
			Name:           c.String("name"),
			Protocol:       pt,
			ProtocolPort:   c.Int("port"),
			LoadBalancerID: c.String("loadbalancer-id"),
			SecretID:       c.String("secret-id"),
			SNISecretID:    c.StringSlice("sni-secret-id"),
			AllowedCIDRS:	c.StringSlice("allowed-cidrs") ,
		}
		if c.IsSet("timeout-client-data") {
			timeoutClientData := c.Int("timeout-client-data")
			opts.TimeoutClientData = &timeoutClientData
		}
		if c.IsSet("timeout-member-connect") {
			timeoutMemberConnect := c.Int("timeout-member-connect")
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}
		if c.IsSet("timeout-member-data") {
			timeoutMemberData := c.Int("timeout-member-data")
			opts.TimeoutMemberData = &timeoutMemberData
		}
		if c.IsSet("connection-limit") {
			connectionLimit := c.Int("connection-limit")
			opts.ConnectionLimit = &connectionLimit
		}
	
		results, err := listeners.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			listenerID, err := listeners.ExtractListenerIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve listener ID from task info: %w", err)
			}
			listener, err := listeners.Get(client, listenerID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get listener with ID: %s. Error: %w", listenerID, err)
			}
			utils.ShowResults(listener, c.String("format"))
			return nil, nil
		})
	},
}

var listenerGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "show loadbalancer listener",
	ArgsUsage: "<listener_id>",
	Category:  "listener",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, listenerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewLBListenerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := listeners.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var listenerDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete loadbalancer listener",
	ArgsUsage: "<listener_id>",
	Category:  "listener",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		listenerID, err := flags.GetFirstStringArg(c, listenerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLBListenerClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := listeners.Delete(client, listenerID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			listener, err := listeners.Get(client, listenerID).Extract()
			if err == nil {
				if listener != nil && listener.IsDeleted() {
					return nil, nil
				}
				return nil, fmt.Errorf("cannot delete listener with ID: %s", listenerID)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})
	},
}

var listenerUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "update loadbalancer listener",
	ArgsUsage: "<listener_id>",
	Category:  "listener",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "listener name",
		},
		&cli.StringFlag{
			Name:    "secret-id",
			Aliases: []string{"s"},
			Usage:   "ID of the secret where PKCS12 file is stored for TERMINATED_HTTPS load balancer",
		},
		&cli.StringSliceFlag{
			Name:  "sni-secret-id",
			Usage: "List of secret's ID containing PKCS12 format certificate/key bundles for TERMINATED_HTTPS listeners",
		},
		&cli.StringSliceFlag{
			Name:  "allowed-cidrs",
			Usage: "List of networks from which listener is accessible",
		},
		&cli.IntFlag{
			Name:     "timeout-client-data",
			Aliases:  []string{"tcd"},
			Usage:    "Frontend client inactivity timeout in milliseconds",
		},
		&cli.IntFlag{
			Name:     "timeout-member-connect",
			Aliases:  []string{"tmc"},
			Usage:    "Backend member connection timeout in milliseconds",
		},
		&cli.IntFlag{
			Name:     "timeout-member-data",
			Aliases:  []string{"tmd"},
			Usage:    "Backend member inactivity timeout in milliseconds",
		},
		&cli.IntFlag{
			Name:     "connection-limit",
			Aliases:  []string{"cl"},
			Usage:    "Limit of the simultaneous connections",
		},
	},
	Action: func(c *cli.Context) error {
		listenerID, err := flags.GetFirstStringArg(c, listenerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewLBListenerClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := listeners.UpdateOpts{
			Name: c.String("name"),
			SecretID: c.String("secret-id"),
			SNISecretID: c.StringSlice("sni-secret-id"),
			AllowedCIDRS: c.StringSlice("allowed-cidrs"),
		}
		if c.IsSet("timeout-client-data") {
			timeoutClientData := c.Int("timeout-client-data")
			opts.TimeoutClientData = &timeoutClientData
		}
		if c.IsSet("timeout-member-connect") {
			timeoutMemberConnect := c.Int("timeout-member-connect")
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}
		if c.IsSet("timeout-member-data") {
			timeoutMemberData := c.Int("timeout-member-data")
			opts.TimeoutMemberData = &timeoutMemberData
		}
		if c.IsSet("connection-limit") {
			connectionLimit := c.Int("connection-limit")
			opts.ConnectionLimit = &connectionLimit
		}

		result, err := listeners.Update(client, listenerID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if result == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ListenerCommands = cli.Command{
	Name:  "listener",
	Usage: "GCloud loadbalancer listeners API",
	Subcommands: []*cli.Command{
		&listenerListSubCommand,
		&listenerGetSubCommand,
		&listenerUpdateSubCommand,
		&listenerDeleteSubCommand,
		&listenerCreateSubCommand,
	},
}
