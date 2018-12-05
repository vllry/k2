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

	err = updateContainerTargetList(db, w.Id, w.Replicas)
	if err != nil {
		return err
	}
	return nil
}

func saveWorkloadSpec(db *dbConn, w *WorkloadSpec) error {
	bytes, err := yaml.Marshal(&w)
	if err != nil {
		return err
	}
	_, err = db.set("/workloads/"+w.Id+"/spec", string(bytes))
	return err
}

// TODO replace this approach. Need more situational awareness with unscheduling.
func generateContainerTargetList(containerList []string, prefix string, replicas int) ([]string, bool) {
	numContainers := len(containerList)
	if numContainers == replicas {
		return containerList, false
	}
	if numContainers > replicas {
		containerList = containerList[:replicas]
	} else if numContainers < replicas {
		for i := 0; i < replicas-numContainers; i++ {
			containerList = append(containerList, prefix+"-"+randStringBytes(8))
		}
	}
	return containerList, true
}

func updateContainerTargetList(db *dbConn, workloadId string, replicas int) error {
	containers, err := getWorkloadContainerList(db, workloadId)
	if err != nil {
		return err
	}
	updatedContainers, changes := generateContainerTargetList(containers, workloadId, replicas)
	if changes {
		err = saveWorkloadContainerList(db, workloadId, updatedContainers)
		return err
	}
	return nil
}

func saveWorkloadContainerList(db *dbConn, workloadId string, containerIds []string) error {
	_, err := db.setYaml("/workloads/"+workloadId+"/containers", containerIds)
	return err
}

func getWorkloadContainerList(db *dbConn, workloadId string) ([]string, error) {
	var containers []string
	found, err := db.getStruct("/workloads/"+workloadId+"/containers", containers)
	if err != nil {
		return nil, err
	}
	if !found {
		return make([]string, 0), nil
	}
	return containers, nil
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
