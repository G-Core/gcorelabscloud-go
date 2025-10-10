package clusters

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/dbaas/postgres/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	taskclient "github.com/G-Core/gcorelabscloud-go/client/tasks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/dbaas/postgres/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	"github.com/urfave/cli/v2"
)

const defaultPGConfigSettings = `
huge_pages=off
max_connections=100
shared_buffers=256MB
effective_cache_size=768MB
maintenance_work_mem=64MB
work_mem=2MB
checkpoint_completion_target=0.9
wal_buffers=-1
min_wal_size=1GB
max_wal_size=4GB
random_page_cost=1.2
effective_io_concurrency=200
`

var defaultRoleAttributes = []clusters.RoleAttribute{
	clusters.RoleAttributeLogin,
	clusters.RoleAttributeInherit,
}

var clusterShowCommand = cli.Command{
	Name:      "show",
	Usage:     "Show details of a PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name>`,
	Action:    showClusterAction,
}

var clusterListCommand = cli.Command{
	Name:     "list",
	Usage:    "List all PostgreSQL clusters",
	Category: "clusters",
	Action:   listClustersAction,
	Flags:    flags.OffsetLimitFlags,
}

var clusterDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete a PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name>`,
	Action:    deleteClusterAction,
	Flags:     flags.WaitCommandFlags,
}

var clusterCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create a new PostgreSQL cluster",
	Category: "clusters",
	Action:   createClusterAction,
	Flags:    createClusterFlags(),
}

var clusterUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update an existing PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name>`,
	Action:    updateClusterAction,
	Flags:     updateClusterFlags(),
}

var clusterAddUserCommand = cli.Command{
	Name:      "add-user",
	Usage:     "Add a user to an existing PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name>`,
	Action:    addUserAction,
	Flags:     addUserFlags(),
}

var clusterRemoveUserCommand = cli.Command{
	Name:      "remove-user",
	Usage:     "Remove a user from an existing PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name> <user_name>`,
	Action:    removeUserAction,
	Flags:     flags.WaitCommandFlags,
}

var clusterAddDatabaseCommand = cli.Command{
	Name:      "add-database",
	Usage:     "Add a database to an existing PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name>`,
	Action:    addDatabaseAction,
	Flags:     addDatabaseFlags(),
}

var clusterRemoveDatabaseCommand = cli.Command{
	Name:      "remove-database",
	Usage:     "Remove a database from an existing PostgreSQL cluster",
	Category:  "clusters",
	ArgsUsage: `<cluster_name> <database_name>`,
	Action:    removeDatabaseAction,
	Flags:     flags.WaitCommandFlags,
}

func showClusterAction(c *cli.Context) error {
	clusterName := c.Args().First()
	if clusterName == "" {
		_ = cli.ShowCommandHelp(c, "show")
		return cli.Exit("cluster_name is required", 1)
	}

	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	clusterDetails := clusters.Get(pgClient, clusterName)
	if clusterDetails.Err != nil {
		return cli.Exit(clusterDetails.Err, 1)
	}

	utils.ShowResults(clusterDetails.Body, c.String("format"))
	return nil
}

func listClustersAction(c *cli.Context) error {
	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	var clustersList []clusters.PostgresSQLClusterShort
	if !c.IsSet("limit") && !c.IsSet("offset") {
		clustersList, err = clusters.ListAll(pgClient, nil)
		if err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		listOpts := clusters.ListOpts{
			Limit:  c.Int("limit"),
			Offset: c.Int("offset"),
		}
		pager := clusters.List(pgClient, listOpts)
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			list, err2 := clusters.ExtractClusters(page)
			if err2 != nil {
				return false, err2
			}
			clustersList = list
			return false, nil // stop after first page
		})
		if err != nil {
			return cli.Exit(err, 1)
		}
	}
	utils.ShowResults(clustersList, c.String("format"))
	return nil
}

func deleteClusterAction(c *cli.Context) error {
	clusterName := c.Args().First()
	if clusterName == "" {
		_ = cli.ShowCommandHelp(c, "delete")
		return cli.Exit("cluster_name is required", 1)
	}
	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	result := clusters.Delete(pgClient, clusterName, clusters.DeleteOpts{})
	if result.Err != nil {
		return cli.Exit(result.Err, 1)
	}
	taskResults, err := result.Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	return utils.WaitTaskAndShowResult(c, pgClient, taskResults, false, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(pgClient, clusterName).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete postgres cluster with name: %s", clusterName)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
}

func createClusterAction(c *cli.Context) error {
	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := clusters.CreateOpts{
		ClusterName: c.String("name"),
		Storage: clusters.PGStorageConfigurationOpts{
			SizeGiB: c.Int("storage-size"),
			Type:    c.String("storage-type"),
		},
		Databases: []clusters.DatabaseOpts{
			{
				Name:  c.String("db-name"),
				Owner: c.String("user-name"),
			},
		},
		PGServerConfiguration: clusters.PGServerConfigurationOpts{
			PGConf:  c.String("pg-config-conf"),
			Version: c.String("pg-config-version"),
		},
		Network: clusters.NetworkOpts{
			ACL:         c.StringSlice("network-acl"),
			NetworkType: "public",
		},
		Flavor: clusters.FlavorOpts{
			CPU:       c.Int("flavor-cpu"),
			MemoryGiB: c.Int("flavor-memory"),
		},
	}

	// add user with custom role attributes if specified, otherwise use default
	roleAttributes := make([]clusters.RoleAttribute, 0)
	if c.IsSet("user-role-attribute") {
		for _, attr := range c.StringSlice("user-role-attribute") {
			if !clusters.IsValidRoleAttribute(attr) {
				return cli.Exit(fmt.Sprintf("invalid role attribute: %s", attr), 1)
			}
			roleAttributes = append(roleAttributes, clusters.RoleAttribute(attr))
		}
	} else {
		roleAttributes = defaultRoleAttributes
	}
	opts.Users = []clusters.PgUserOpts{
		{
			Name:           c.String("user-name"),
			RoleAttributes: roleAttributes,
		},
	}

	if c.IsSet("pg-config-pooler-mode") {
		opts.PGServerConfiguration.Pooler = &clusters.PoolerOpts{
			Mode: clusters.PoolerMode(c.String("pg-config-pooler-mode")),
			Type: "pgbouncer",
		}
	}

	if c.IsSet("ha-replication-mode") {
		opts.HighAvailability = &clusters.HighAvailabilityOpts{
			ReplicationMode: clusters.HighAvailabilityReplicationMode(c.String("ha-replication-mode")),
		}
	}

	result := clusters.Create(pgClient, opts)
	if err := handleTaskResult(c, result, pgClient, opts.ClusterName); err != nil {
		return err
	}
	return nil
}

func updateClusterAction(c *cli.Context) error {
	clusterName := c.Args().First()
	if clusterName == "" {
		_ = cli.ShowCommandHelp(c, "update")
		return cli.Exit("cluster_name is required", 1)
	}

	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	var cluster *clusters.PostgresSQLCluster
	if (c.IsSet("db-name") && c.IsSet("db-owner")) ||
		(c.IsSet("user-name") && c.IsSet("user-role-attribute")) {
		cluster, err = clusters.Get(pgClient, clusterName).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
	}

	opts := clusters.UpdateOpts{}

	if c.IsSet("storage-size") {
		opts.Storage = &clusters.PGStorageConfigurationUpdateOpts{
			SizeGiB: c.Int("storage-size"),
		}
	}
	// change the db owner
	if c.IsSet("db-name") && c.IsSet("db-owner") {
		if cluster == nil {
			return cli.Exit("internal error: cluster is not fetched", 1)
		}
		foundOwner := false
		for _, user := range cluster.Users {
			if user.Name == c.String("db-owner") {
				foundOwner = true
			}
		}
		if !foundOwner {
			return cli.Exit(fmt.Sprintf("db-owner with name '%s' not found in cluster '%s'",
				c.String("db-owner"), clusterName), 1)
		}

		newDatabases := make([]clusters.DatabaseOpts, 0)
		foundDatabase := false
		for _, db := range cluster.Databases {
			dbOpts := clusters.DatabaseOpts{Name: db.Name}
			if db.Name == c.String("db-name") {
				if db.Owner == c.String("db-owner") {
					return cli.Exit(fmt.Sprintf("database '%s' in cluster '%s' already has owner '%s'",
						c.String("db-name"), clusterName, c.String("db-owner")), 1)
				}
				dbOpts.Owner = c.String("db-owner")
				foundDatabase = true
			} else {
				dbOpts.Owner = db.Owner
			}
			newDatabases = append(newDatabases, dbOpts)
		}
		if !foundDatabase {
			return cli.Exit(fmt.Sprintf("database with name '%s' not found in cluster '%s'",
				c.String("db-name"), clusterName), 1)
		}
		opts.Databases = newDatabases
	}
	// change the user role attributes
	if c.IsSet("user-name") && c.IsSet("user-role-attribute") {
		if cluster == nil {
			return cli.Exit("internal error: cluster is not fetched", 1)
		}
		roleAttributes := make([]clusters.RoleAttribute, 0)
		for _, attr := range c.StringSlice("user-role-attribute") {
			if !clusters.IsValidRoleAttribute(attr) {
				return cli.Exit(fmt.Sprintf("invalid role attribute: %s", attr), 1)
			}
			roleAttributes = append(roleAttributes, clusters.RoleAttribute(attr))
		}
		newUsers := make([]clusters.PgUserOpts, 0)
		foundUser := false
		for _, user := range cluster.Users {
			userOpts := clusters.PgUserOpts{
				Name:           user.Name,
				RoleAttributes: user.RoleAttributes,
			}
			if user.Name == c.String("user-name") {
				foundUser = true
				userOpts.RoleAttributes = roleAttributes
			}
			newUsers = append(newUsers, userOpts)
		}
		if !foundUser {
			return cli.Exit(fmt.Sprintf("user with name '%s' not found in cluster '%s'",
				c.String("user-name"), clusterName), 1)
		}
		opts.Users = newUsers
	}
	if c.IsSet("pg-config-pooler-mode") {
		opts.PGServerConfiguration.Pooler = &clusters.PoolerOpts{
			Mode: clusters.PoolerMode(c.String("pg-config-pooler-mode")),
			Type: "pgbouncer",
		}
	}
	if c.IsSet("ha-replication-mode") {
		opts.HighAvailability = &clusters.HighAvailabilityOpts{
			ReplicationMode: clusters.HighAvailabilityReplicationMode(c.String("ha-replication-mode")),
		}
	}
	if c.IsSet("flavor-cpu") && c.IsSet("flavor-memory") {
		opts.Flavor = &clusters.FlavorOpts{
			CPU:       c.Int("flavor-cpu"),
			MemoryGiB: c.Int("flavor-memory"),
		}
	}
	if c.IsSet("network-acl") {
		opts.Network = &clusters.NetworkOpts{
			ACL:         c.StringSlice("network-acl"),
			NetworkType: "public",
		}
	}

	result := clusters.Update(pgClient, clusterName, opts)
	if err := handleTaskResult(c, result, pgClient, clusterName); err != nil {
		return err
	}
	return nil
}

func addUserAction(c *cli.Context) error {
	clusterName := c.Args().First()
	if clusterName == "" {
		_ = cli.ShowCommandHelp(c, "add-user")
		return cli.Exit("cluster_name is required", 1)
	}

	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// first, get the existing users
	cluster, err := clusters.Get(pgClient, clusterName).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	newUsers := make([]clusters.PgUserOpts, 0)
	for _, user := range cluster.Users {
		if user.Name == c.String("user-name") {
			return cli.Exit(fmt.Sprintf("user with name %s already exists in cluster %s",
				c.String("user-name"), clusterName), 1)
		}
		newUsers = append(newUsers, clusters.PgUserOpts{
			Name:           user.Name,
			RoleAttributes: user.RoleAttributes,
		})
	}

	// second, add the new user
	newUser := clusters.PgUserOpts{Name: c.String("user-name")}
	if c.IsSet("user-role-attribute") {
		newUser.RoleAttributes = make([]clusters.RoleAttribute, 0)
		for _, attr := range c.StringSlice("user-role-attribute") {
			if !clusters.IsValidRoleAttribute(attr) {
				return cli.Exit(fmt.Sprintf("invalid role attribute: %s", attr), 1)
			}
			newUser.RoleAttributes = append(newUser.RoleAttributes, clusters.RoleAttribute(attr))
		}
	} else {
		newUser.RoleAttributes = defaultRoleAttributes
	}
	newUsers = append(newUsers, newUser)

	// finally, update the cluster with the new list of users
	result := clusters.Update(pgClient, clusterName, clusters.UpdateOpts{Users: newUsers})
	if err := handleTaskResult(c, result, pgClient, clusterName); err != nil {
		return err
	}
	return nil
}

func removeUserAction(c *cli.Context) error {
	clusterName := c.Args().Get(0)
	userName := c.Args().Get(1)
	if clusterName == "" || userName == "" {
		_ = cli.ShowCommandHelp(c, "remove-user")
		return cli.Exit("cluster_name and user_name are required", 1)
	}

	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// first, get the existing users
	cluster, err := clusters.Get(pgClient, clusterName).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	newUsers := make([]clusters.PgUserOpts, 0)
	for _, user := range cluster.Users {
		if user.Name != userName {
			newUsers = append(newUsers, clusters.PgUserOpts{
				Name:           user.Name,
				RoleAttributes: user.RoleAttributes,
			})
		}
	}
	if len(newUsers) == len(cluster.Users) {
		return cli.Exit(fmt.Sprintf("user with name '%s' does not exist in cluster '%s'",
			userName, clusterName), 1)
	}

	// finally, update the cluster with the new list of users
	result := clusters.Update(pgClient, clusterName, clusters.UpdateOpts{Users: newUsers})
	if err := handleTaskResult(c, result, pgClient, clusterName); err != nil {
		return err
	}
	return nil
}

func addDatabaseAction(c *cli.Context) error {
	clusterName := c.Args().First()
	if clusterName == "" {
		_ = cli.ShowCommandHelp(c, "add-database")
		return cli.Exit("cluster_name is required", 1)
	}

	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// first, get the existing databases
	cluster, err := clusters.Get(pgClient, clusterName).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	newDatabases := make([]clusters.DatabaseOpts, 0)
	for _, db := range cluster.Databases {
		if db.Name == c.String("db-name") {
			return cli.Exit(fmt.Sprintf("database with name '%s' already exists in cluster '%s'",
				c.String("db-name"), clusterName), 1)
		}
		newDatabases = append(newDatabases, clusters.DatabaseOpts{
			Name:  db.Name,
			Owner: db.Owner,
		})
	}

	// check that the owner exists
	foundOwner := false
	for _, user := range cluster.Users {
		if user.Name == c.String("db-owner") {
			foundOwner = true
		}
	}
	if !foundOwner {
		return cli.Exit(fmt.Sprintf("db-owner with name '%s' not found in cluster '%s'",
			c.String("db-owner"), clusterName), 1)
	}

	// second, add the new database
	newDatabases = append(newDatabases, clusters.DatabaseOpts{
		Name:  c.String("db-name"),
		Owner: c.String("db-owner"),
	})
	// finally, update the cluster with the new list of databases
	result := clusters.Update(pgClient, clusterName, clusters.UpdateOpts{Databases: newDatabases})
	if err := handleTaskResult(c, result, pgClient, clusterName); err != nil {
		return err
	}
	return nil
}

func removeDatabaseAction(c *cli.Context) error {
	clusterName := c.Args().Get(0)
	databaseName := c.Args().Get(1)
	if clusterName == "" || databaseName == "" {
		_ = cli.ShowCommandHelp(c, "remove-database")
		return cli.Exit("cluster_name and database_name are required", 1)
	}

	pgClient, err := client.NewPostgresClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// first, get the existing databases
	cluster, err := clusters.Get(pgClient, clusterName).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	newDatabases := make([]clusters.DatabaseOpts, 0)
	foundDatabase := false
	for _, db := range cluster.Databases {
		if db.Name != databaseName {
			newDatabases = append(newDatabases, clusters.DatabaseOpts{
				Name:  db.Name,
				Owner: db.Owner,
			})
		} else {
			foundDatabase = true
		}
	}
	if !foundDatabase {
		return cli.Exit(fmt.Sprintf("database with name '%s' does not exist in cluster '%s'",
			databaseName, clusterName), 1)
	}

	// finally, update the cluster with the new list of databases
	result := clusters.Update(pgClient, clusterName, clusters.UpdateOpts{Databases: newDatabases})
	if err := handleTaskResult(c, result, pgClient, clusterName); err != nil {
		return err
	}
	return nil
}

func handleTaskResult(c *cli.Context, result tasks.Result, pgClient *gcorecloud.ServiceClient, clusterName string) error {
	if result.Err != nil {
		return cli.Exit(result.Err, 1)
	}
	taskResults, err := result.Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}
	return utils.WaitTaskAndShowResult(c, taskClient, taskResults, true,
		waitForClusterOperation(pgClient, clusterName))
}

func waitForClusterOperation(pgClient *gcorecloud.ServiceClient, clusterName string) func(task tasks.TaskID) (interface{}, error) {
	return func(task tasks.TaskID) (interface{}, error) {
		cluster, err := clusters.Get(pgClient, clusterName).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot perform operation on PostgreSQL cluster with name: %s. Error: %w",
				clusterName, err)
		}
		return cluster, nil
	}
}

func createClusterFlags() []cli.Flag {
	createFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "name of the PostgreSQL cluster",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "storage-size",
			Usage:    "size of the storage in GB (min: 1, max: 100)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "storage-type",
			Usage:    "type of the storage (e.g. ssd-hiiops, standard, ssd-lowlatency)",
			Required: false,
			Value:    "ssd-hiiops",
		},
		&cli.StringFlag{
			Name:     "db-name",
			Usage:    "name of the initial database",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "user-name",
			Usage:    "name of the initial user and owner of the database",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "user-role-attribute",
			Usage:    "role attributes for the initial user (can be repeated)",
			Required: false,
			Value:    cli.NewStringSlice(clusters.RoleAttributeSliceToStrings(defaultRoleAttributes)...),
		},
		&cli.StringFlag{
			Name:        "pg-config-conf",
			Usage:       `pg.conf settings (e.g. "\nlisten_addresses = 'localhost'\nport = 5432")`,
			Required:    false,
			Value:       defaultPGConfigSettings,
			DefaultText: "the default PostgreSQL configuration settings",
		},
		&cli.StringFlag{
			Name:     "pg-config-version",
			Usage:    "version of PostgreSQL (e.g. 13, 14, 15)",
			Required: false,
			Value:    "15",
		},
		&cli.StringFlag{
			Name:     "pg-config-pooler-mode",
			Usage:    "mode of the connection pooler (e.g. session, statement, transaction)",
			Value:    "session",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "network-acl",
			Usage:    "list of IP addresses or CIDR blocks to allow access to the cluster (can be repeated)",
			Required: false,
			Value:    cli.NewStringSlice("0.0.0.0/0"),
		},
		&cli.IntFlag{
			Name:     "flavor-cpu",
			Usage:    "number of CPU cores (min: 1)",
			Required: false,
			Value:    1,
		},
		&cli.IntFlag{
			Name:     "flavor-memory",
			Usage:    "amount of RAM in GiB (min: 1)",
			Required: false,
			Value:    1,
		},
		&cli.StringFlag{
			Name:     "ha-replication-mode",
			Usage:    "high availability replication mode (e.g. async, sync)",
			Required: false,
		},
	}
	createFlags = append(createFlags, flags.WaitCommandFlags...)
	return createFlags
}

func updateClusterFlags() []cli.Flag {
	updateFlags := []cli.Flag{
		&cli.IntFlag{
			Name:        "storage-size",
			Usage:       "size of the storage in GB (min: 1, max: 100). Can only be increased.",
			Required:    false,
			DefaultText: "not set",
		},
		&cli.StringFlag{
			Name:     "db-name",
			Usage:    "name of the database to update the owner for",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "db-owner",
			Usage:    "new owner of the database specified with --db-name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-name",
			Usage:    "name of the user for which to update role attributes",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "user-role-attribute",
			Usage:    "role attributes for the user specified with --user-name (can be repeated)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-config-pooler-mode",
			Usage:    "mode of the connection pooler (e.g. session, statement, transaction)",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "network-acl",
			Usage:    "list of IP addresses or CIDR blocks to allow access to the cluster (can be repeated)",
			Required: false,
		},
		&cli.IntFlag{
			Name:        "flavor-cpu",
			Usage:       "number of CPU cores (min: 1)",
			Required:    false,
			DefaultText: "not set",
		},
		&cli.IntFlag{
			Name:        "flavor-memory",
			Usage:       "amount of RAM in GiB (min: 1)",
			Required:    false,
			DefaultText: "not set",
		},
		&cli.StringFlag{
			Name:     "ha-replication-mode",
			Usage:    "high availability replication mode (e.g. async, sync)",
			Required: false,
		},
	}
	updateFlags = append(updateFlags, flags.WaitCommandFlags...)
	return updateFlags
}

func addUserFlags() []cli.Flag {
	addFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "user-name",
			Usage:    "name of the user to add",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "user-role-attribute",
			Usage:    "role attributes for the initial user (can be repeated)",
			Required: false,
			Value:    cli.NewStringSlice(clusters.RoleAttributeSliceToStrings(defaultRoleAttributes)...),
		},
	}
	addFlags = append(addFlags, flags.WaitCommandFlags...)
	return addFlags
}

func addDatabaseFlags() []cli.Flag {
	addFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "db-name",
			Usage:    "name of the database to add",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "db-owner",
			Usage:    "owner of the database to add (must be an existing user)",
			Required: true,
		},
	}
	addFlags = append(addFlags, flags.WaitCommandFlags...)
	return addFlags
}

var Commands = cli.Command{
	Name:        "clusters",
	Usage:       "Manage PostgreSQL clusters",
	Description: "Commands for managing PostgreSQL clusters",
	Subcommands: []*cli.Command{
		&clusterShowCommand,
		&clusterListCommand,
		&clusterDeleteCommand,
		&clusterCreateCommand,
		&clusterUpdateCommand,
		&clusterAddUserCommand,
		&clusterRemoveUserCommand,
		&clusterAddDatabaseCommand,
		&clusterRemoveDatabaseCommand,
	},
}
