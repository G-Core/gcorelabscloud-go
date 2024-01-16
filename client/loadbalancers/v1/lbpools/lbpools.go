package lbpools

import (
	"fmt"
	"net"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var (
	lbpoolIDText           = "pool_id is mandatory argument"
	memberIDText           = "member_id is mandatory argument"
	protocolTypes          = types.ProtocolType("").StringList()
	loadBalancerAlgorithms = types.LoadBalancerAlgorithm("").StringList()
	healthMonitorTypes     = types.HealthMonitorType("").StringList()
	httpMethods            = types.HTTPMethod("").StringList()
	persistenceTypes       = types.PersistenceType("").StringList()
)

func getHealthMonitor(c *cli.Context) (*lbpools.CreateHealthMonitorOpts, error) {

	healthMonitorType, err := types.HealthMonitorType(c.String("healthmonitor-type")).ValidOrNil()
	if err != nil || healthMonitorType == nil {
		return nil, err
	}

	healthMonitorDelay := c.Int("healthmonitor-delay")
	if healthMonitorDelay == 0 {
		return nil, fmt.Errorf("--healthmonitor-delay should be set for health monitor %s", healthMonitorType)
	}
	healthMonitorMaxRetires := c.Int("healthmonitor-max-retries")
	if healthMonitorMaxRetires == 0 {
		return nil, fmt.Errorf("--healthmonitor-max-retries should be set for health monitor %s", healthMonitorType)
	}
	healthMonitorTimeout := c.Int("healthmonitor-timeout")
	if healthMonitorTimeout == 0 {
		return nil, fmt.Errorf("--healthmonitor-timeout should be set for health monitor %s", healthMonitorType)
	}
	hm := lbpools.CreateHealthMonitorOpts{
		Type:           *healthMonitorType,
		Delay:          healthMonitorDelay,
		MaxRetries:     healthMonitorMaxRetires,
		Timeout:        healthMonitorTimeout,
		MaxRetriesDown: c.Int("healthmonitor-max-retries-down"),
		ExpectedCodes:  c.String("healthmonitor-expected-codes"),
	}
	if healthMonitorType.IsHTTPType() {
		httpMethod := types.HTTPMethod(c.String("healthmonitor-http-method"))
		if err := httpMethod.IsValid(); err != nil {
			return nil, err
		}
		hm.HTTPMethod = types.HTTPMethodPointer(httpMethod)
		httpMethodURLPath := c.String("healthmonitor-url-path")
		if httpMethodURLPath == "" {
			return nil, fmt.Errorf("--healthmonitor-url-path should be set for health monitor type %s", healthMonitorType)
		}
		hm.URLPath = httpMethodURLPath
	}
	return &hm, nil
}

func getSessionPersistence(c *cli.Context) (*lbpools.CreateSessionPersistenceOpts, error) {

	sessionPersistenceType, err := types.PersistenceType(c.String("session-persistence-type")).ValidOrNil()
	if err != nil || sessionPersistenceType == nil {
		return nil, err
	}

	sessionPersistenceCookiesName := c.String("session-cookies-name")
	if sessionPersistenceType.ISCookiesType() && sessionPersistenceCookiesName == "" {
		return nil, fmt.Errorf("--session-cookies-name should be set for session persistence type %s", sessionPersistenceType)
	}

	return &lbpools.CreateSessionPersistenceOpts{
		PersistenceGranularity: c.String("session-persistence-granularity"),
		PersistenceTimeout:     c.Int("session-persistence-timeout"),
		Type:                   *sessionPersistenceType,
		CookieName:             sessionPersistenceCookiesName,
	}, nil
}

func getPoolMembers(c *cli.Context) ([]lbpools.CreatePoolMemberOpts, error) {
	memberAddresses := c.StringSlice("member-address")
	if len(memberAddresses) == 0 {
		return nil, nil
	}
	memberPorts := c.IntSlice("member-port")
	if len(memberAddresses) != len(memberPorts) {
		return nil, fmt.Errorf("number of --member-address should be equal --member-port")
	}
	memberWeights := c.IntSlice("member-weight")
	memberSubnetIDs := c.StringSlice("member-subnet-id")
	memberInstanceIDs := c.StringSlice("member-instance-id")
	var members []lbpools.CreatePoolMemberOpts

	type addressPortPair struct {
		ip   string
		port int
	}

	mp := map[addressPortPair]int{}

	for idx, addr := range memberAddresses {
		memberAddr := net.ParseIP(addr)
		if memberAddr == nil {
			return nil, fmt.Errorf("malformed member-address %s", addr)
		}
		member := lbpools.CreatePoolMemberOpts{
			Address:      memberAddr,
			ProtocolPort: memberPorts[idx],
			Weight: func(idx int) int {
				if idx < len(memberWeights) {
					return memberWeights[idx]
				}
				return 0
			}(idx),
			InstanceID: func(idx int) string {
				if idx < len(memberInstanceIDs) {
					return memberInstanceIDs[idx]
				}
				return ""
			}(idx),
			SubnetID: func(idx int) string {
				if idx < len(memberSubnetIDs) {
					return memberInstanceIDs[idx]
				}
				return ""
			}(idx),
		}
		members = append(members, member)
		mp[addressPortPair{
			ip:   addr,
			port: memberPorts[idx],
		}]++
	}

	for key, value := range mp {
		if value > 1 {
			return nil, fmt.Errorf("same address and port %s:%d have been set %d times", key.ip, key.port, value)
		}
	}

	return members, nil

}

var lbpoolListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "loadbalancer pools list",
	Category: "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "loadbalancer-id",
			Aliases:     []string{"l"},
			Usage:       "loadbalancer ID",
			Required:    false,
			DefaultText: "<nil>",
		},
		&cli.StringFlag{
			Name:        "listener-id",
			Usage:       "listener ID",
			Required:    false,
			DefaultText: "<nil>",
		},
		&cli.BoolFlag{
			Name:        "details",
			Usage:       "show details",
			Required:    false,
			DefaultText: "<nil>",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := lbpools.ListOpts{
			LoadBalancerID: utils.StringToPointer(c.String("loadbalancer-id")),
			ListenerID:     utils.StringToPointer(c.String("listener-id")),
			MemberDetails:  utils.BoolToPointer(c.Bool("details")),
		}

		results, err := lbpools.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var lbpoolGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "show loadbalancer pool",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := lbpools.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var lbpoolDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete loadbalancer pool",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		lbpoolID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := lbpools.Delete(client, lbpoolID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			lbpool, err := lbpools.Get(client, lbpoolID).Extract()
			if err == nil {
				if lbpool != nil && lbpool.IsDeleted() {
					return nil, nil
				}
				return nil, fmt.Errorf("cannot delete lbpool with ID: %s", lbpoolID)
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

var lbpoolCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "create loadbalancer pool",
	Category: "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "lbpool name",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "protocol",
			Aliases: []string{"p"},
			Value: &utils.EnumValue{
				Enum:    protocolTypes,
				Default: types.ProtocolTypeTCP.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypes, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "algorithm",
			Aliases: []string{"a"},
			Value: &utils.EnumValue{
				Enum:    loadBalancerAlgorithms,
				Default: types.LoadBalancerAlgorithmRoundRobin.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(loadBalancerAlgorithms, ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "loadbalancer",
			Aliases:  []string{"lb"},
			Usage:    "loadbalancer ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "listener",
			Aliases:  []string{"lbl"},
			Usage:    "loadbalancer listener ID",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-type",
			Aliases: []string{"hmt"},
			Value: &utils.EnumValue{
				Enum: healthMonitorTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(healthMonitorTypes, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-delay",
			Aliases:  []string{"hmd"},
			Usage:    "health monitor checking delay",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries",
			Aliases:  []string{"hmr"},
			Usage:    "health monitor checking max retries",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-timeout",
			Aliases:  []string{"hmto"},
			Usage:    "health monitor checking timeout",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries-down",
			Aliases:  []string{"hmrd"},
			Usage:    "health monitor checking max retries down",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-http-method",
			Aliases: []string{"hmhm"},
			Value: &utils.EnumValue{
				Enum:    httpMethods,
				Default: types.HTTPMethodGET.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-url-path",
			Aliases:  []string{"hmup"},
			Usage:    "health monitor checking url path",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-expected-codes",
			Aliases:  []string{"hmec"},
			Usage:    "health monitor checking expected codes",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "session-persistence-type",
			Aliases: []string{"spt"},
			Value: &utils.EnumValue{
				Enum: persistenceTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(persistenceTypes, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-cookies-name",
			Aliases:  []string{"scn"},
			Usage:    "health monitor session persistence cookies name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "session-persistence-timeout",
			Aliases:  []string{"spto"},
			Usage:    "health monitor session persistence timeout",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-persistence-granularity",
			Aliases:  []string{"spg"},
			Usage:    "health monitor session persistence granularity",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-address",
			Aliases:  []string{"ma"},
			Usage:    "pool member address",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-port",
			Aliases:  []string{"mp"},
			Usage:    "pool member port",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-weight",
			Aliases:  []string{"mw"},
			Usage:    "pool member weight",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-subnet-id",
			Aliases:  []string{"ms"},
			Usage:    "pool subnet ID",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-instance-id",
			Aliases:  []string{"mi"},
			Usage:    "pool instance ID",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-client-data",
			Aliases:  []string{"tcd"},
			Usage:    "frontend client inactivity timeout in milliseconds",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-member-connect",
			Aliases:  []string{"tmc"},
			Usage:    "backend member connection timeout in milliseconds",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-member-data",
			Aliases:  []string{"tmd"},
			Usage:    "backend member inactivity timeout in milliseconds",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		pt := types.ProtocolType(c.String("protocol"))
		if err := pt.IsValid(); err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		lba := types.LoadBalancerAlgorithm(c.String("algorithm"))
		if err := lba.IsValid(); err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		members, err := getPoolMembers(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		hm, err := getHealthMonitor(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		sp, err := getSessionPersistence(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		if members == nil {
			members = []lbpools.CreatePoolMemberOpts{}
		}

		loadBalancerID := utils.StringToPointer(c.String("loadbalancer"))
		listenerID := utils.StringToPointer(c.String("listener"))

		if loadBalancerID == nil && listenerID == nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(fmt.Errorf("either --loadbalancer or --listener should be set"), 1)
		}

		timeoutClientData := c.Int("timeout-client-data")
		timeoutMemberConnect := c.Int("timeout-member-connect")
		timeoutMemberData := c.Int("timeout-member-data")

		opts := lbpools.CreateOpts{
			Name:                 c.String("name"),
			Protocol:             pt,
			LBPoolAlgorithm:      lba,
			Members:              members,
			LoadBalancerID:       c.String("loadbalancer"),
			ListenerID:           c.String("listener"),
			HealthMonitor:        hm,
			SessionPersistence:   sp,
			TimeoutClientData:    timeoutClientData,
			TimeoutMemberConnect: timeoutMemberConnect,
			TimeoutMemberData:    timeoutMemberData,
		}

		results, err := lbpools.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			lbpoolID, err := lbpools.ExtractPoolIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve lbpool ID from task info: %w", err)
			}
			lbpool, err := lbpools.Get(client, lbpoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", lbpoolID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var lbpoolCreateMemberSubCommand = cli.Command{
	Name:      "create",
	Usage:     "create loadbalancer pool member",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "address",
			Aliases:  []string{"a"},
			Usage:    "pool member address",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "port",
			Aliases:  []string{"p"},
			Usage:    "pool member port",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "weight",
			Aliases:  []string{"mw"},
			Usage:    "pool member weight",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "subnet-id",
			Aliases:  []string{"s"},
			Usage:    "pool subnet ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "instance-id",
			Aliases:  []string{"i"},
			Usage:    "pool instance ID",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		lbpoolID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		address := net.ParseIP(c.String("address"))
		if address == nil {
			return cli.NewExitError(fmt.Errorf("malformed address %s", c.String("address")), 1)
		}

		opts := lbpools.CreatePoolMemberOpts{
			Address:      address,
			ProtocolPort: c.Int("port"),
			Weight:       c.Int("weight"),
			SubnetID:     c.String("subnet-id"),
			InstanceID:   c.String("instance-id"),
		}

		results, err := lbpools.CreateMember(client, lbpoolID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			memberID, err := lbpools.ExtractPoolMemberIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve lbpool member ID from task info: %w", err)
			}
			lbpool, err := lbpools.Get(client, lbpoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", memberID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var lbpoolDeleteMemberSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete loadbalancer pool",
	ArgsUsage: "<member_id>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "pool-id",
			Aliases:  []string{"p"},
			Usage:    "pool ID",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		memberID, err := flags.GetFirstStringArg(c, memberIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		lbpoolID := c.String("pool-id")
		results, err := lbpools.DeleteMember(client, lbpoolID, memberID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			lbpool, err := lbpools.Get(client, lbpoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get loadbalancer pool with ID: %s", lbpoolID)
			}
			members := lbpool.Members
			for _, m := range members {
				if m.ID == memberID {
					return nil, fmt.Errorf("cannot delete loadbalancer pool member with ID: %s", memberID)
				}
			}
			return nil, nil
		})
	},
}

var lbpoolUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "update loadbalancer pool",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "lbpool name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-client-data",
			Aliases:  []string{"tcd"},
			Usage:    "frontend client inactivity timeout in milliseconds",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-member-connect",
			Aliases:  []string{"tmc"},
			Usage:    "backend member connection timeout in milliseconds",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout-member-data",
			Aliases:  []string{"tmd"},
			Usage:    "backend member inactivity timeout in milliseconds",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "algorithm",
			Aliases: []string{"a"},
			Value: &utils.EnumValue{
				Enum: loadBalancerAlgorithms,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(loadBalancerAlgorithms, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-type",
			Aliases: []string{"hmt"},
			Value: &utils.EnumValue{
				Enum: healthMonitorTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(healthMonitorTypes, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-delay",
			Aliases:  []string{"hmd"},
			Usage:    "health monitor checking delay",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries",
			Aliases:  []string{"hmr"},
			Usage:    "health monitor checking max retries",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-timeout",
			Aliases:  []string{"hmto"},
			Usage:    "health monitor checking timeout",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries-down",
			Aliases:  []string{"hmrd"},
			Usage:    "health monitor checking max retries down",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-http-method",
			Aliases: []string{"hmhm"},
			Value: &utils.EnumValue{
				Enum:    httpMethods,
				Default: types.HTTPMethodGET.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-expected-codes",
			Aliases:  []string{"hmec"},
			Usage:    "health monitor checking expected codes",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-url-path",
			Aliases:  []string{"hmup"},
			Usage:    "health monitor checking url path",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "session-persistence-type",
			Aliases: []string{"spt"},
			Value: &utils.EnumValue{
				Enum: persistenceTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-cookies-name",
			Aliases:  []string{"scn"},
			Usage:    "health monitor session persistence cookies name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "session-persistence-timeout",
			Aliases:  []string{"spto"},
			Usage:    "health monitor session persistence timeout",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-persistence-granularity",
			Aliases:  []string{"spg"},
			Usage:    "health monitor session persistence granularity",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-address",
			Aliases:  []string{"ma"},
			Usage:    "pool member address",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-port",
			Aliases:  []string{"mp"},
			Usage:    "pool member port",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-weight",
			Aliases:  []string{"mw"},
			Usage:    "pool member weight",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-instance-id",
			Aliases:  []string{"mi"},
			Usage:    "pool instance ID",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		lbPoolID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		lba, err := types.LoadBalancerAlgorithm(c.String("algorithm")).ValidOrNil()
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		members, err := getPoolMembers(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		hm, err := getHealthMonitor(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		sp, err := getSessionPersistence(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		timeoutClientData := c.Int("timeout-client-data")
		timeoutMemberConnect := c.Int("timeout-member-connect")
		timeoutMemberData := c.Int("timeout-member-data")

		opts := lbpools.UpdateOpts{
			Name:                 c.String("name"),
			Members:              members,
			HealthMonitor:        hm,
			SessionPersistence:   sp,
			TimeoutClientData:    timeoutClientData,
			TimeoutMemberConnect: timeoutMemberConnect,
			TimeoutMemberData:    timeoutMemberData,
		}

		if lba != nil {
			opts.LBPoolAlgorithm = *lba
		}

		results, err := lbpools.Update(client, lbPoolID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			_, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			lbpool, err := lbpools.Get(client, lbPoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", lbPoolID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var lbpoolUnsetSubCommand = cli.Command{
	Name:      "unset",
	Usage:     "unset loadbalancer pool fields",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:     "session-persistence",
			Aliases:  []string{"sp"},
			Usage:    "Disables session persistence on the pool",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		lbPoolID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "unset")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := lbpools.UnsetOpts{
			SessionPersistence: c.Bool("session-persistence"),
		}
		results, err := lbpools.Unset(client, lbPoolID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			_, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			lbpool, err := lbpools.Get(client, lbPoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", lbPoolID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var lbpoolCreateHealthMonitorSubCommand = cli.Command{
	Name:      "create",
	Usage:     "create pool's health monitor",
	ArgsUsage: "<pool_id>",
	Category:  "healthmonitor",
	Flags: append([]cli.Flag{
		&cli.GenericFlag{
			Name:    "healthmonitor-type",
			Aliases: []string{"hmt"},
			Value: &utils.EnumValue{
				Enum: healthMonitorTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(healthMonitorTypes, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-delay",
			Aliases:  []string{"hmd"},
			Usage:    "health monitor checking delay",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries",
			Aliases:  []string{"hmr"},
			Usage:    "health monitor checking max retries",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-timeout",
			Aliases:  []string{"hmto"},
			Usage:    "health monitor checking timeout",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries-down",
			Aliases:  []string{"hmrd"},
			Usage:    "health monitor checking max retries down",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-http-method",
			Aliases: []string{"hmhm"},
			Value: &utils.EnumValue{
				Enum:    httpMethods,
				Default: types.HTTPMethodGET.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-expected-codes",
			Aliases:  []string{"hmec"},
			Usage:    "health monitor checking expected codes",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-url-path",
			Aliases:  []string{"hmup"},
			Usage:    "health monitor checking url path",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		lbPoolID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		hm, err := getHealthMonitor(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		results, err := lbpools.CreateHealthMonitor(client, lbPoolID, hm).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			_, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			lbpool, err := lbpools.Get(client, lbPoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", lbPoolID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var lbpoolDeleteHealthMonitorSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete loadbalancer pool's health monitor",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		lbpoolID, err := flags.GetFirstStringArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLBPoolClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		if err = lbpools.DeleteHealthMonitor(client, lbpoolID).ExtractErr(); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var PoolCommands = cli.Command{
	Name:  "pool",
	Usage: "GCloud loadbalancer pools API",
	Subcommands: []*cli.Command{
		&lbpoolListSubCommand,
		&lbpoolGetSubCommand,
		&lbpoolUpdateSubCommand,
		&lbpoolDeleteSubCommand,
		&lbpoolCreateSubCommand,
		&lbpoolUnsetSubCommand,
		{
			Name:  "member",
			Usage: "GCloud loadbalancer pool members API",
			Subcommands: []*cli.Command{
				&lbpoolCreateMemberSubCommand,
				&lbpoolDeleteMemberSubCommand,
			},
		},
		{
			Name:  "healthmonitor",
			Usage: "GCloud pool's health monitor API",
			Subcommands: []*cli.Command{
				&lbpoolCreateHealthMonitorSubCommand,
				&lbpoolDeleteHealthMonitorSubCommand,
			},
		},
	},
}
