package utils

import "github.com/se2022-qiaqia/course-system/model/req"

type group struct {
	startWeek uint
	endWeek   uint
	schedules []*req.CourseSchedule
}

// group 记录了周范围尽可能广的课程安排组若干个，对于每个组来说，
// 虽然包含的课程安排时间对应的周次不完全相同，但是必须至少有一个课程安排在同一周。、
type groups struct {
	groups []*group
}

func newGroups() groups {
	return groups{}
}

func (g *groups) findGroup(startWeek, endWeek uint) (int, bool) {
	for i, g2 := range g.groups {
		// 无交集的情况，再取反，得到就是有交集的情况
		if !(g2.endWeek < startWeek || g2.startWeek > endWeek) {
			return i, true
		}
	}
	return -1, false
}

func (g *groups) join(schedule *req.CourseSchedule) {
	startWeek := schedule.StartWeekId
	endWeek := schedule.EndWeekId
	index, ok := g.findGroup(startWeek, endWeek)
	if !ok {
		g.groups = append(g.groups, &group{
			startWeek: startWeek,
			endWeek:   endWeek,
			schedules: []*req.CourseSchedule{schedule},
		})
		return
	}

	g.groups[index].schedules = append(g.groups[index].schedules, schedule)
	// 尽可能拓宽组
	if g.groups[index].endWeek < endWeek {
		g.groups[index].endWeek = endWeek
	}
	if g.groups[index].startWeek > startWeek {
		g.groups[index].startWeek = startWeek
	}
}

func IsScheduleConflict(schedules []*req.CourseSchedule, ignoreWeek bool) bool {
	if ignoreWeek {
		// memo: 7 * 24 * 4 = 672 bytes
		days := make([][]bool, 7)
		for i := 0; i < 7; i++ {
			days[i] = make([]bool, 24+1) // req.CourseSchedule 里的HoursId从1开始，这里额外开辟一个空间来处理，而不对HoursId进行偏移
		}

		for _, schedule := range schedules {
			for i := schedule.StartHoursId; i <= schedule.EndHoursId; i++ {
				if days[schedule.DayOfWeek][i] {
					return true
				}
				days[schedule.DayOfWeek][i] = true
			}
		}
	} else {
		gs := newGroups()
		for _, schedule := range schedules {
			gs.join(schedule)
		}
		for _, g := range gs.groups {
			if IsScheduleConflict(g.schedules, true) {
				return true
			}
		}
	}
	return false
}
