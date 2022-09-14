package predicates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"

	kubevirtv1 "kubevirt.io/api/core/v1"

	"github.com/project-flotta/osbuild-operator/internal/predicates"
)

var _ = Describe("OSBuildEnvConfig VM reconciliation predicate", func() {
	DescribeTable("should reconcile update", func(old, new runtimeclient.Object) {
		// given
		p := predicates.OSBuildEnvConfigJobFinished{}
		e := event.UpdateEvent{ObjectOld: old, ObjectNew: new}

		// when
		shouldReconcile := p.Update(e)

		// then
		Expect(shouldReconcile).To(BeTrue())
	},
		Entry("when the conditions changed",
			&kubevirtv1.VirtualMachine{
				Status: kubevirtv1.VirtualMachineStatus{
					Conditions: []kubevirtv1.VirtualMachineCondition{
						{
							Type: kubevirtv1.VirtualMachineReady,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
			&kubevirtv1.VirtualMachine{
				Status: kubevirtv1.VirtualMachineStatus{
					Conditions: []kubevirtv1.VirtualMachineCondition{
						{
							Type: kubevirtv1.VirtualMachineReady,
							Status: v1.ConditionTrue,
						},
					},
				},
			},
		),
	)

	DescribeTable("should not reconcile update", func(old, new runtimeclient.Object) {
		// given
		p := predicates.OSBuildEnvConfigJobFinished{}
		e := event.UpdateEvent{ObjectOld: old, ObjectNew: new}

		// when
		shouldReconcile := p.Update(e)

		// then
		Expect(shouldReconcile).To(BeFalse())
	},
		Entry("when conditions are the exactly the same",
			&kubevirtv1.VirtualMachine{
				Status: kubevirtv1.VirtualMachineStatus{
					Conditions: []kubevirtv1.VirtualMachineCondition{
						{
							Type: kubevirtv1.VirtualMachineReady,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
			&kubevirtv1.VirtualMachine{
				Status: kubevirtv1.VirtualMachineStatus{
					Conditions: []kubevirtv1.VirtualMachineCondition{
						{
							Type: kubevirtv1.VirtualMachineReady,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
		),
		Entry("when only the order of the conditions changed",
			&kubevirtv1.VirtualMachine{
				Status: kubevirtv1.VirtualMachineStatus{
					Conditions: []kubevirtv1.VirtualMachineCondition{
						{
							Type: kubevirtv1.VirtualMachineReady,
							Status: v1.ConditionFalse,
						},
						{
							Type: kubevirtv1.VirtualMachinePaused,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
			&kubevirtv1.VirtualMachine{
				Status: kubevirtv1.VirtualMachineStatus{
					Conditions: []kubevirtv1.VirtualMachineCondition{
						{
							Type: kubevirtv1.VirtualMachinePaused,
							Status: v1.ConditionFalse,
						},
						{
							Type: kubevirtv1.VirtualMachineReady,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
		),
		Entry("when old is missing",
			nil,
			&kubevirtv1.VirtualMachine{},
		),
		Entry("when new is missing",
			&kubevirtv1.VirtualMachine{},
			nil,
		),
		Entry("when new is not Job",
			&kubevirtv1.VirtualMachine
			&batchv1.CronJob{},
		),
	)

	It("should not reconcile create", func() {
		// given
		p := predicates.OSBuildEnvConfigJobFinished{}
		e := event.CreateEvent{
			Object: &batchv1.Job{
				Status: batchv1.JobStatus{
					Active: 1,
				},
			},
		}

		// when
		shouldReconcile := p.Create(e)

		// then
		Expect(shouldReconcile).To(BeFalse())
	})

	It("should not reconcile delete", func() {
		// given
		p := predicates.OSBuildEnvConfigJobFinished{}
		e := event.DeleteEvent{
			Object: &batchv1.Job{
				Status: batchv1.JobStatus{
					Active: 1,
				},
			},
		}

		// when
		shouldReconcile := p.Delete(e)

		// then
		Expect(shouldReconcile).To(BeFalse())
	})

	It("should not reconcile for generic event", func() {
		// given
		// given
		p := predicates.OSBuildEnvConfigJobFinished{}
		e := event.GenericEvent{
			Object: &batchv1.Job{
				Status: batchv1.JobStatus{
					Active: 1,
				},
			},
		}

		// when
		shouldReconcile := p.Generic(e)

		// then
		Expect(shouldReconcile).To(BeFalse())
	})
})
