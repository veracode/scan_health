package report

import "strings"

func getReportModule(r *Report, moduleName string) *Module {
	for index, reportModule := range r.Modules {
		if strings.EqualFold(moduleName, reportModule.Name) {
			return &r.Modules[index]
		}
	}

	// Module has not been found so add it

	module := Module{
		Name: moduleName,
	}

	r.Modules = append(r.Modules, module)

	for index, reportModule := range r.Modules {
		if strings.EqualFold(moduleName, reportModule.Name) {
			return &r.Modules[index]
		}
	}

	return nil
}

func (r *Report) AddModuleInstance(moduleName string, moduleInstance ModuleInstance) {
	module := getReportModule(r, moduleName)
	module.Instances = append(module.Instances, moduleInstance)
}
