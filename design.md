# Cron Horizontal Pod Autoscaler

Table of Contents
=================

   * [Cron Horizontal Pod Autoscaler](#cron-horizontal-pod-autoscaler)
      * [Summary](#summary)
      * [Motivation](#motivation)
         * [Goals](#goals)
         * [Non-Goals](#non-goals)
      * [Proposal](#proposal)
         * [User Stories](#user-stories)
            * [Story 1](#story-1)
            * [Story 2](#story-2)
         * [Implementation Details/Notes/Constraints](#implementation-detailsnotesconstraints)
      * [Future work](#future-work)

## Summary

Cron Horizontal Pod Autoscaler(CronHPA) enables us to auto scale workloads(those support `scale`subresource, e.g. deployment, statefulset) periodically using [crontab](https://en.wikipedia.org/wiki/Cron) scheme.

## Motivation

A lot of game players will play games from Friday evening to Sunday evening. It will be better to provide a better experience for players if game servers(pods) could be scaled to a lager size at Friday evening and scaled to normal size at Sunday evening. And that's what game server admins do every week.

Some other customers also do similar things because their products' usage also have peaks and valleys periodically. Time based auto scheduler could scale pods in advance, and will provide a better experience for users. CronHPA is for their requirement.

### Goals

* Users could auto scale their workloads those support `scale` subresource periodically using crontab scheme.

### Non-Goals

* CronHPA does not support workloads that does not support `scale` subresource now.

## Proposal

At first, we'd like to use [CronJob](https://gs.io/docs/concepts/workloads/controllers/cron-jobs/) to implement this feature. Howerver `CronJob` only supports one `schedule` and users must specify an image and commands, it will be a little inappropriate for above use cases.

We propose a new CRD `CronHPA` which is inspired by `CronJob` (thanks a lot!) to meet users' requirement. Users could specify as many `schedule`s as they want.

`CronHPA` example:

```
apiVersion: gs.io/v1
kind: CronHPA
metadata:
  name: example-cron-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: demo-deployment
  crons:
    - schedule: "0 23 * * 5"  // Set replicas to 60 every Friday 23:00
      targetReplicas: 60
    - schedule: "0 23 * * 7"  // Set replicas to 30 every Sunday 23:00
      targetReplicas: 30
```

### User Stories

#### Story 1

As a game server admin, I find there will be a lot of game players logining to our servers and playing games from Friday 20:00 to Sunday 23:00 every week, I'd better to auto scale game servers at that time. And I will use following `CronHPA` to solve my problem.

```
apiVersion: gs.io/v1
kind: CronHPA
metadata:
  name: game-servers-cronhpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: game-servers
  crons:
    - schedule: "0 20 * * 5"  // Set replicas to 60 every Friday 20:00
      targetReplicas: 60
    - schedule: "0 23 * * 7"  // Set replicas to 30 every Sunday 23:00
      targetReplicas: 30
```

#### Story 2

As a web server admin, I find pageviews for my web is high at every day 8:00~9:00 am and 7:00~9:00 pm. And I will use following `CronHPA` to solve my problem.

```
apiVersion: gs.io/v1
kind: CronHPA
metadata:
  name: web-servers-cronhpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web-servers
  crons:
    - schedule: "0 8 * * *"  // Set replicas to 60 every day 8:00
      targetReplicas: 60
    - schedule: "0 9 * * *"  // Set replicas to 10 every day 9:00
      targetReplicas: 10
    - schedule: "0 19 * * *"  // Set replicas to 60 every day 19:00
      targetReplicas: 60
    - schedule: "0 21 * * *"  // Set replicas to 10 every day 21:00
      targetReplicas: 10
```

### Implementation Details/Notes/Constraints

We have implenented the corresponding controller manager which is similar to `CroJob` controller manager.

CRD `CronHPA` realated data structure is defined as following:

```
// CronHPA represents a set of crontabs to set target's replicas.
type CronHPA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired identities of pods in this cronhpa.
	Spec CronHPASpec `json:"spec,omitempty"`

	// Status is the current status of pods in this CronHPA. This data
	// may be out of date by some window of time.
	Status CronHPAStatus `json:"status,omitempty"`
}

// A CronHPASpec is the specification of a CronHPA.
type CronHPASpec struct {
	// scaleTargetRef points to the target resource to scale
	ScaleTargetRef autoscalingv2.CrossVersionObjectReference `json:"scaleTargetRef" protobuf:"bytes,1,opt,name=scaleTargetRef"`

	Crons []Cron `json:"crons" protobuf:"bytes,2,opt,name=crons"`
}

type Cron struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule" protobuf:"bytes,1,opt,name=schedule"`

	TargetReplicas int32 `json:"targetReplicas" protobuf:"varint,2,opt,name=targetReplicas"`
}

// CronHPAStatus represents the current state of a CronHPA.
type CronHPAStatus struct {
	// Information when was the last time the schedule was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty" protobuf:"bytes,2,opt,name=lastScheduleTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPAList is a collection of CronHPA.
type CronHPAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronHPA `json:"items"`
}
```

In the conroller, it will list `CronHPA` periodically, and scale workloads using `scale` subresource if needed.

## Future work

We will try to integrate `HPA` and `CronHPA` for a better experience.
