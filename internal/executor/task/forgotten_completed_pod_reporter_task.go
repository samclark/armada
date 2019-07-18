package task

import (
	"errors"
	"github.com/G-Research/k8s-batch/internal/executor/domain"
	"github.com/G-Research/k8s-batch/internal/executor/reporter"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	listers "k8s.io/client-go/listers/core/v1"
	"time"
)

type ForgottenCompletedPodReporterTask struct {
	PodLister     listers.PodLister
	EventReporter reporter.EventReporter
}

func (podCleanup ForgottenCompletedPodReporterTask) Run() {
	podsToBeReported := getAllPodsRequiringCompletionEvent(podCleanup.PodLister)

	for _, pod := range podsToBeReported {
		podCleanup.EventReporter.ReportEvent(pod)
	}
}

func getAllPodsRequiringCompletionEvent(podLister listers.PodLister) []*v1.Pod {
	requirement, err := labels.NewRequirement(domain.JobId, selection.Exists, []string{})
	if err != nil {
		return nil
		//TODO Handle error case
	}

	selector := labels.NewSelector().Add(*requirement)
	allBatchPodsNotMarkedForCleanup, err := podLister.List(selector)

	if err != nil {
		//TODO Do something in case of error
	}

	completedBatchPodsNotMarkedForCleanup := filterCompletedPods(allBatchPodsNotMarkedForCleanup)
	completedBatchPodsToBeReported := filterPodsInStateForLongerThanGivenDuration(completedBatchPodsNotMarkedForCleanup, time.Minute*2)

	return completedBatchPodsToBeReported
}

func filterCompletedPods(pods []*v1.Pod) []*v1.Pod {
	completedPods := make([]*v1.Pod, 0, len(pods))

	for _, pod := range pods {
		if isInTerminalState(pod) {
			completedPods = append(completedPods, pod)
		}
	}

	return completedPods
}

func filterPodsInStateForLongerThanGivenDuration(pods []*v1.Pod, duration time.Duration) []*v1.Pod {
	podsInStateForLongerThanDuration := make([]*v1.Pod, 0)

	expiryTime := time.Now().Add(-duration)
	for _, pod := range pods {
		lastStatusChange, err := lastStatusChange(pod)
		if err != nil || lastStatusChange.Before(expiryTime) {
			podsInStateForLongerThanDuration = append(podsInStateForLongerThanDuration, pod)
		}
	}

	return podsInStateForLongerThanDuration
}

func lastStatusChange(pod *v1.Pod) (time.Time, error) {
	conditions := pod.Status.Conditions

	if len(conditions) <= 0 {
		return *new(time.Time), errors.New("no state changes found, cannot determine last status change")
	}

	var maxStatusChange time.Time

	for _, condition := range conditions {
		if condition.LastTransitionTime.Time.After(maxStatusChange) {
			maxStatusChange = condition.LastTransitionTime.Time
		}
	}

	return maxStatusChange, nil
}