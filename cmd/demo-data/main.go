package main

import (
	"fmt"
	"github.com/se2022-qiaqia/course-system/config"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/services"
	"strconv"
	"strings"
)

func main() {
	config.Init()
	log.Init()
	dao.Init()
	dao.Migrate()

	fmt.Printf("您的操作将对数据库产生影响, 请确保您的数据库已备份\n")
	if config.Config.Database.Sqlite != nil {
		fmt.Printf("数据库类型：%s\n", "sqlite")
		fmt.Printf("数据库文件：%s\n", config.Config.Database.Sqlite.Filename)
	} else {
		fmt.Println("未知数据库类型")
	}
	fmt.Printf("您确定要执行吗？(y/N)")

	var input string
	fmt.Scanln(&input)
	if strings.ToLower(input) != "y" {
		fmt.Printf("操作取消！\n")
		return
	}

	collegeAnonymous := dao.College{
		Model: dao.Model{
			ID: 1,
		},
		Name: "<未分配学院>",
	}
	collegeMath := &dao.College{Name: "数学院"}
	collegePhysical := &dao.College{Name: "物理学院"}
	collegeChemistry := &dao.College{Name: "化学院"}
	collegeBiology := &dao.College{Name: "生物学院"}
	collegeComputerScience := &dao.College{Name: "计算机科学学院"}
	collegeSoftwareEngineering := &dao.College{Name: "软件学院"}
	collegeElectronicEngineering := &dao.College{Name: "电子工程学院"}
	collegeBatch := []*dao.College{collegeMath, collegePhysical, collegeChemistry, collegeBiology, collegeComputerScience, collegeSoftwareEngineering, collegeElectronicEngineering}
	dao.DB.FirstOrCreate(&collegeAnonymous, "id = ?", collegeAnonymous.ID)
	dao.DB.Create(collegeBatch)

	adminDefault := dao.User{
		Model: dao.Model{
			ID: 1,
		},
		Username:     "admin",
		RealName:     "关丽媛",
		College:      collegeAnonymous,
		EntranceYear: 2000,
	}
	adminDefault.SetPassword("hello-admin")
	dao.DB.FirstOrCreate(&adminDefault, "id = ?", adminDefault.ID)
	teacherMath1 := &dao.User{
		Username:     "math-teacher1",
		RealName:     "数学教师1",
		Role:         dao.RoleTeacher,
		College:      *collegeMath,
		EntranceYear: 2019,
	}
	teacherMath2 := &dao.User{
		Username:     "math-teacher2",
		RealName:     "数学教师2",
		Role:         dao.RoleTeacher,
		College:      *collegeMath,
		EntranceYear: 2019,
	}
	teacherPhysical1 := &dao.User{
		Username:     "physical-teacher1",
		RealName:     "物理教师1",
		Role:         dao.RoleTeacher,
		College:      *collegePhysical,
		EntranceYear: 2019,
	}
	teacherPhysical2 := &dao.User{
		Username:     "physical-teacher2",
		RealName:     "物理教师2",
		Role:         dao.RoleTeacher,
		College:      *collegePhysical,
		EntranceYear: 2019,
	}
	teacherChemistry1 := &dao.User{
		Username:     "chemistry-teacher1",
		RealName:     "化学教师1",
		Role:         dao.RoleTeacher,
		College:      *collegeChemistry,
		EntranceYear: 2019,
	}
	teacherChemistry2 := &dao.User{
		Username:     "chemistry-teacher2",
		RealName:     "化学教师2",
		Role:         dao.RoleTeacher,
		College:      *collegeChemistry,
		EntranceYear: 2019,
	}
	teacherBiology1 := &dao.User{
		Username:     "biology-teacher1",
		RealName:     "生物教师1",
		Role:         dao.RoleTeacher,
		College:      *collegeBiology,
		EntranceYear: 2019,
	}
	teacherBiology2 := &dao.User{
		Username:     "biology-teacher2",
		RealName:     "生物教师2",
		Role:         dao.RoleTeacher,
		College:      *collegeBiology,
		EntranceYear: 2019,
	}
	teacherComputerScience1 := &dao.User{
		Username:     "computer-science-teacher1",
		RealName:     "计算机科学教师1",
		Role:         dao.RoleTeacher,
		College:      *collegeComputerScience,
		EntranceYear: 2019,
	}
	teacherComputerScience2 := &dao.User{
		Username:     "computer-science-teacher2",
		RealName:     "计算机科学教师2",
		Role:         dao.RoleTeacher,
		College:      *collegeComputerScience,
		EntranceYear: 2019,
	}
	teacherSoftwareEngineering1 := &dao.User{
		Username:     "software-engineering-teacher1",
		RealName:     "软件工程教师1",
		Role:         dao.RoleTeacher,
		College:      *collegeSoftwareEngineering,
		EntranceYear: 2019,
	}
	teacherSoftwareEngineering2 := &dao.User{
		Username:     "software-engineering-teacher2",
		RealName:     "软件工程教师2",
		Role:         dao.RoleTeacher,
		College:      *collegeSoftwareEngineering,
		EntranceYear: 2019,
	}
	teacherElectronicEngineering1 := &dao.User{
		Username:     "electronic-engineering-teacher1",
		RealName:     "电子工程教师1",
		Role:         dao.RoleTeacher,
		College:      *collegeElectronicEngineering,
		EntranceYear: 2019,
	}
	teacherElectronicEngineering2 := &dao.User{
		Username:     "electronic-engineering-teacher2",
		RealName:     "电子工程教师2",
		Role:         dao.RoleTeacher,
		College:      *collegeElectronicEngineering,
		EntranceYear: 2019,
	}
	teacherBatch := []*dao.User{
		teacherMath1, teacherMath2,
		teacherPhysical1, teacherPhysical2,
		teacherChemistry1, teacherChemistry2,
		teacherBiology1, teacherBiology2,
		teacherComputerScience1, teacherComputerScience2,
		teacherSoftwareEngineering1, teacherSoftwareEngineering2,
		teacherElectronicEngineering1, teacherElectronicEngineering2,
	}
	for _, batch := range teacherBatch {
		batch.SetPassword("demo-password")
	}
	dao.DB.Create(teacherBatch)

	ccAdvancedMath := &dao.CourseCommon{
		Name:    "高等数学",
		Credits: 5,
		Hours:   60,
		College: *collegeMath,
	}
	ccLinearAlgebra := &dao.CourseCommon{
		Name:    "线性代数",
		Credits: 3,
		Hours:   60,
		College: *collegeMath,
	}
	ccProbabilityTheory := &dao.CourseCommon{
		Name:    "概率论与数理统计",
		Credits: 3,
		Hours:   60,
		College: *collegeMath,
	}
	ccUniversityPhysics := &dao.CourseCommon{
		Name:    "大学物理",
		Credits: 5,
		Hours:   60,
		College: *collegePhysical,
	}
	ccQuantumPhysics := &dao.CourseCommon{
		Name:    "量子物理",
		Credits: 3,
		Hours:   60,
		College: *collegePhysical,
	}
	ccOrganicChemistry := &dao.CourseCommon{
		Name:    "有机化学",
		Credits: 3,
		Hours:   60,
		College: *collegeChemistry,
	}
	ccInorganicChemistry := &dao.CourseCommon{
		Name:    "无机化学",
		Credits: 3,
		Hours:   60,
		College: *collegeChemistry,
	}
	ccUniversityBiology := &dao.CourseCommon{
		Name:    "大学生物",
		Credits: 2,
		Hours:   60,
		College: *collegeBiology,
	}
	ccComputerScienceIntro := &dao.CourseCommon{
		Name:    "计算机科学导论",
		Credits: 3,
		Hours:   48,
		College: *collegeComputerScience,
	}
	ccEmbeddedSystem := &dao.CourseCommon{
		Name:    "嵌入式系统",
		Credits: 3.5,
		Hours:   48,
		College: *collegeComputerScience,
	}
	ccDataStructureForCS := &dao.CourseCommon{
		Name:    "数据结构",
		Credits: 3.5,
		Hours:   48,
		College: *collegeComputerScience,
	}
	ccDataStructureForSE := &dao.CourseCommon{
		Name:    "数据结构",
		Credits: 3,
		Hours:   48,
		College: *collegeSoftwareEngineering,
	}
	ccJavaEE := &dao.CourseCommon{
		Name:    "JavaEE企业级应用开发",
		Credits: 3,
		Hours:   48,
		College: *collegeSoftwareEngineering,
	}
	ccCircuitAnalysis := &dao.CourseCommon{
		Name:    "电路分析",
		Credits: 3,
		Hours:   48,
		College: *collegeElectronicEngineering,
	}
	ccDigitalElectronics := &dao.CourseCommon{
		Name:    "数字电子技术",
		Credits: 3,
		Hours:   48,
		College: *collegeElectronicEngineering,
	}
	ccBatch := []*dao.CourseCommon{
		ccAdvancedMath, ccLinearAlgebra, ccProbabilityTheory,
		ccUniversityPhysics, ccQuantumPhysics, ccOrganicChemistry,
		ccInorganicChemistry, ccUniversityBiology, ccComputerScienceIntro,
		ccEmbeddedSystem, ccDataStructureForCS, ccDataStructureForSE,
		ccJavaEE, ccCircuitAnalysis, ccDigitalElectronics,
	}
	dao.DB.Create(ccBatch)

	semester := &dao.Semester{
		Year: 2022,
		Term: 1,
	}
	dao.DB.Create(semester)
	services.Services.Setting.Set(services.KeyCurrentSemester, strconv.Itoa(int(semester.ID)))

	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccAdvancedMath,
		Teacher:      *teacherMath1,
		Location:     "A105",
		Quota:        80,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   3,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
			{
				DayOfWeek:   5,
				HoursId:     1,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccAdvancedMath,
		Teacher:      *teacherMath2,
		Location:     "A106",
		Quota:        80,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   3,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
			{
				DayOfWeek:   5,
				HoursId:     1,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccAdvancedMath,
		Teacher:      *teacherMath2,
		Location:     "A106",
		Quota:        80,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
			{
				DayOfWeek:   4,
				HoursId:     1,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccLinearAlgebra,
		Teacher:      *teacherMath1,
		Location:     "A105",
		Quota:        80,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
			{
				DayOfWeek:   4,
				HoursId:     1,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccLinearAlgebra,
		Teacher:      *teacherMath2,
		Location:     "A105",
		Quota:        80,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     1,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccProbabilityTheory,
		Teacher:      *teacherMath1,
		Location:     "A105",
		Quota:        80,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     6,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccUniversityPhysics,
		Teacher:      *teacherPhysical1,
		Location:     "B105",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     6,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccUniversityPhysics,
		Teacher:      *teacherPhysical2,
		Location:     "B105",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     1,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccQuantumPhysics,
		Teacher:      *teacherPhysical2,
		Location:     "B106",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   1,
				HoursId:     1,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccQuantumPhysics,
		Teacher:      *teacherPhysical1,
		Location:     "B104",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   1,
				HoursId:     3,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccOrganicChemistry,
		Teacher:      *teacherChemistry1,
		Location:     "A204",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   1,
				HoursId:     6,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccOrganicChemistry,
		Teacher:      *teacherChemistry2,
		Location:     "A203",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     6,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccInorganicChemistry,
		Teacher:      *teacherChemistry1,
		Location:     "A204",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   1,
				HoursId:     1,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccInorganicChemistry,
		Teacher:      *teacherChemistry2,
		Location:     "A203",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   2,
				HoursId:     1,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccUniversityBiology,
		Teacher:      *teacherBiology1,
		Location:     "A207",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   5,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccComputerScienceIntro,
		Teacher:      *teacherComputerScience1,
		Location:     "A208",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   5,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccComputerScienceIntro,
		Teacher:      *teacherComputerScience2,
		Location:     "A209",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   5,
				HoursId:     3,
				HoursCount:  2,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccDataStructureForCS,
		Teacher:      *teacherComputerScience1,
		Location:     "A208",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   5,
				HoursId:     5,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccDataStructureForCS,
		Teacher:      *teacherComputerScience2,
		Location:     "A208",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   4,
				HoursId:     5,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccEmbeddedSystem,
		Teacher:      *teacherComputerScience2,
		Location:     "A208",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   3,
				HoursId:     1,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccJavaEE,
		Teacher:      *teacherSoftwareEngineering1,
		Location:     "A310",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   3,
				HoursId:     6,
				HoursCount:  5,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccJavaEE,
		Teacher:      *teacherSoftwareEngineering2,
		Location:     "A311",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   3,
				HoursId:     6,
				HoursCount:  5,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccCircuitAnalysis,
		Teacher:      *teacherElectronicEngineering1,
		Location:     "A302",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   4,
				HoursId:     6,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccCircuitAnalysis,
		Teacher:      *teacherElectronicEngineering2,
		Location:     "A303",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   4,
				HoursId:     6,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccDigitalElectronics,
		Teacher:      *teacherElectronicEngineering1,
		Location:     "A302",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   4,
				HoursId:     11,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})
	dao.DB.Create(&dao.CourseSpecific{
		CourseCommon: *ccDigitalElectronics,
		Teacher:      *teacherElectronicEngineering2,
		Location:     "A303",
		Quota:        40,
		QuotaUsed:    0,
		Semester:     *semester,
		CourseSchedules: []*dao.CourseSchedule{
			{
				DayOfWeek:   4,
				HoursId:     11,
				HoursCount:  3,
				StartWeekId: 1,
				EndWeekId:   16,
			},
		},
	})

}
