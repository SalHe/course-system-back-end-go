package utils

import (
	"github.com/se2022-qiaqia/course-system/model/req"
	"testing"
)

func TestIsScheduleConflict(t *testing.T) {
	// type args struct {
	// 	schedules  []*req.CourseSchedule
	// 	ignoreWeek bool
	// }
	tests := []struct {
		name string
		// args args
		schedules []*req.CourseSchedule
		conflict  bool
	}{
		{
			name: "不同周=>不冲突",
			schedules: []*req.CourseSchedule{
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
				{
					StartWeekId: 16,
					EndWeekId:   17,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
			},
			conflict: false,
		},
		{
			name: "不同天=>不冲突",
			schedules: []*req.CourseSchedule{
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 2,

					StartHoursId: 1,
					EndHoursId:   5,
				},
			},
			conflict: false,
		},
		{
			name: "同天不同时间段=>不冲突",
			schedules: []*req.CourseSchedule{
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 1, // 1-5
					EndHoursId:   5,
				},
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 6, // 6-10
					EndHoursId:   10,
				},
			},
			conflict: false,
		},

		{
			name: "全同周同天同时间段=>冲突",
			schedules: []*req.CourseSchedule{
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
			},
			conflict: true,
		},
		{
			name: "部分同周、含同天同时间段=>冲突",
			schedules: []*req.CourseSchedule{
				{
					StartWeekId: 1,
					EndWeekId:   15,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
				{
					StartWeekId: 15,
					EndWeekId:   16,

					DayOfWeek: 1,

					StartHoursId: 1,
					EndHoursId:   5,
				},
			},
			conflict: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsScheduleConflict(tt.schedules, false); got != tt.conflict {
				t.Errorf("IsScheduleConflict() = %v, conflict %v", got, tt.conflict)
			}
		})
	}
}
