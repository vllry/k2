package main

func podSetDbPath(id string) string {
	return "/podsets/" + id
}

type PodSetSpec struct {
	Id         string
	Containers []WorkloadContainerSpec
	Replicas   int
}

func newPodSetSpec(w *WorkloadSpec) PodSetSpec {
	return PodSetSpec{
		Id:w.Id + randStringBytes(4),
		Replicas: w.Replicas,
		Containers: w.Containers,
	}
}

func savePodSetSpec(db *dbConn, podset *PodSetSpec) error {
	_, err := db.setYaml(podSetDbPath(podset.Id), podset)
	return err
}
