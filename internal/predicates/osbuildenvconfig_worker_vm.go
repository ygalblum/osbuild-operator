package predicates

import (
	"sort"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/google/go-cmp/cmp"

	kubevirtv1 "kubevirt.io/api/core/v1"
)

type OSBuildEnvConfigWorkerVM struct {
	predicate.Funcs
}

func (OSBuildEnvConfigWorkerVM) Create(e event.CreateEvent) bool {
	return false
}

func (OSBuildEnvConfigWorkerVM) Delete(e event.DeleteEvent) bool {
	return true
}

func (OSBuildEnvConfigWorkerVM) Generic(e event.GenericEvent) bool {
	return false
}

func (OSBuildEnvConfigWorkerVM) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil {
		return false
	}
	if e.ObjectNew == nil {
		return false
	}

	oldVM, ok := e.ObjectOld.(*kubevirtv1.VirtualMachine)
	if !ok {
		return false
	}

	newVM, ok := e.ObjectNew.(*kubevirtv1.VirtualMachine)
	if !ok {
		return false
	}

	unsortedConditions := [][]kubevirtv1.VirtualMachineCondition{
		oldVM.Status.Conditions,
		newVM.Status.Conditions,
	}

	sortedConditions := make([][]kubevirtv1.VirtualMachineCondition, 2)
	for i := range sortedConditions {
		sortedConditions[i] = append([]kubevirtv1.VirtualMachineCondition{}, unsortedConditions[i]...)
		sort.SliceStable(sortedConditions[i], func(j, k int) bool {
			return sortedConditions[i][j].Type < sortedConditions[i][k].Type
		})
	}

	return !cmp.Equal(sortedConditions[0], sortedConditions[1])
}
