# Concepts

## Pod

Directly analogous to Kubernetes - a group of 1+ containers scheduled together.

## PodSet

A PodSet is analogous to a ReplicaSet. It maintains the spec for a group of pods (how many and what parameters each has).

## Workload

A workload is a unit of work that creates PodSets to manage pods, currently analogous to a Kubernetes deployment. The design of workloads is in flux, but they will likely be designed as a base unit, with more opinionated systems (persistent services, batch jobs, etc) implemented on top.

# Controllers

## Workload Controller

The workload controller is responsible for tracking workloads. It creates PodSets to manage individual pods in the deployment.

## PodSet Controller

Manages the spec for individual pods. Given a workload's pod definition and a replica target, the PodSet controller maintains a list of desired pods.

## Schedule Controller

The schedule controller handles the lifecycle of individual pods (and therefore containers), as designed by the PodSet controller.

The controller will create pods that do not exist, and reap pods which are no longer specified to exist.

# Questions

## Workload and PodSet relationship

What's the simplest way to represent different kinds of workloads in a uniform way?

Top level workload concepts:
* services
* stateful
    * stable volume bindings
    * stable network address
* jobs
    * no replication?
    * run to exit 0
* FaaS
    * replace most container spec items

Workloads need the PodSet to handle rolling updates, as a workload may have disparate pod specs at the same time.