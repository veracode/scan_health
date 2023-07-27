package report

import "strings"

func (r *Report) AddModuleInstance(moduleName string, moduleInstance ModuleInstance) {
	for _, reportModule := range r.Modules {
		if strings.EqualFold(moduleName, reportModule.Name) {
			return
		}
	}

	// Module has not been found so add it
	module := Module{
		Name: moduleName,
	}

	module.Instances = append(module.Instances, moduleInstance)

	r.Modules = append(r.Modules, module)
}
