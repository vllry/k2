package main

// Manages workloads for the system.
// Writes container spec information for the scheduler.

// This file is starting to feel like a function overload...

import (
	"gopkg.in/yaml.v2"
)

type WorkloadSpec struct {
	SpecVersionMajor int
	SpecVersionMinor int
	Id               string
	Containers       []WorkloadContainerSpec
	Replicas         int
}

type WorkloadContainerSpec struct {
	Name  string
	Image string
	Tag   string
}

// Creates/updates a workload.
func submitWorkload(db *dbConn, w *WorkloadSpec) error {
	err := saveWorkloadSpec(db, w)
	if err != nil {
		return err
	}
	return nil
}

// Create (or overwrite existing) workload spec.
func saveWorkloadSpec(db *dbConn, w *WorkloadSpec) error {
	bytes, err := yaml.Marshal(&w)
	if err != nil {
		return err
	}
	_, err = db.set("/workloads/"+w.Id+"/spec", string(bytes))
	if err != nil {
		return err
	}

	podset := newPodSetSpec(w)
	err = savePodSetSpec(db, &podset)
	return err
}

//func fetchWorkloadSpec(db *dbConn, workloadId string) (*WorkloadSpec, error) {
//	specYaml, found, err := db.get("/workloads/" + workloadId + "/spec")
//	if err != nil || !found {
//		return nil, err
//	}
//	var spec WorkloadSpec
//	err = yaml.Unmarshal([]byte(specYaml), spec)
//	if err != nil {
//		return nil, err
//	}
//	return &spec, nil
//}
