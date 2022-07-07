package drivers

import "gorm.io/gorm"

type (
	Driver interface {
		Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB
		Name() string
	}

	Factory struct {
		debug   bool
		drivers map[string]Driver
	}
)

func New(debug bool) *Factory {
	factory := Factory{debug: debug, drivers: make(map[string]Driver)}

	factory.Add(Mysql("mysql"))
	factory.Add(PostgreSql("postgresql"))

	return &factory
}

func (d Factory) Register(drivers []Driver) {
	for _, v := range drivers {
		d.Add(v)
	}
}

func (d Factory) Add(driver Driver) {
	d.drivers[driver.Name()] = driver
}

func (d Factory) Connect(driver string, host string, port int, user string, password string, dbname string) *gorm.DB {
	for k, v := range d.drivers {
		if k == driver {
			return v.Connect(host, port, user, password, dbname, d.debug)
		}
	}

	return nil
}
