package report

import (
	"github.com/antfie/scan_health/v2/utils"
	"strings"
)

func getFormattedModuleName(moduleName string) string {
	// Module names derived from flaw modules (from detailed report) can have a version number appended to the exe e.g. abc.exe#1.0.0.0
	if strings.Contains(moduleName, ".exe#") {
		moduleName = strings.Split(moduleName, ".exe#")[0] + ".exe"
	}

	return moduleName
}

func getReportModule(r *Report, moduleName string) *Module {
	moduleName = getFormattedModuleName(moduleName)

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

func (r *Report) AddModuleDependency(moduleName, thisModuleIsDependentOn string) {
	module := getReportModule(r, moduleName)

	dependsOn := getFormattedModuleName(thisModuleIsDependentOn)

	if !utils.IsStringInStringArray(dependsOn, module.DependencyOf) {
		module.DependencyOf = append(module.DependencyOf, dependsOn)
	}
}
