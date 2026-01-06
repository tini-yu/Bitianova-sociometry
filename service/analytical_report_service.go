package service

import (
	"fmt"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/ZeroHawkeye/wordZero/pkg/style"
)

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) SaveReportFile(fullPath string) error {

	doc := document.New()
	// Титульник
	subtitle := doc.AddParagraph("Аналитическая справка по результатам социометрии ")
	subtitle.SetStyle(style.StyleSubtitle)

	//n1
	n1_1 := fmt.Sprintf("Дата проведения диагностики – %s",
		n1.n1_1)
	n1_2 := fmt.Sprintf("Кол-во участников диагностики – %d", n1.n1_2)
	content := doc.AddParagraph(n1_1)
	content = doc.AddParagraph(n1_2)

	content.SetStyle(style.StyleNormal)
	//n2
	content = doc.AddParagraph("Получены следующие эмпирические данные:\n")
	content.SetStyle(style.StyleNormal)

	listn_2_1 := fmt.Sprintf("В классе %.1f%% принятых детей , из них %.1f%% звезд, %.1f%% - популярных.  %.1f%%  непринятых детей, из них:\n",
		n2.n2_1_1, n2.n2_1_2, n2.n2_1_3, n2.n2_1_4)
	content = doc.AddNumberedList(listn_2_1, 0, document.ListTypeDecimal)
	content.SetStyle(style.StyleNormal)

	n2_2 := fmt.Sprintf("%.1f%% отверженных детей, которых воспринимают весьма экспрессивно, но отрицательно, определенное неприятие человека, его качеств, свойств и привычек;", n2.n2_2)
	n2_3 := fmt.Sprintf("%.1f%% изолированных детей, которых нет в эмоциональном реестре группы ни на уровне чувств, ни на уровне отношений;", n2.n2_3)
	n2_4 := fmt.Sprintf("%.1f%% обучающихся непринятых детей из целевых групп, что может свидетельствовать о дефицитах в интеллектуальных, личностных, физических или материальных ресурсов и являться препятствием в расположении одноклассников;", n2.n2_4)
	content = doc.AddBulletList(n2_2, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n2_3, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n2_4, 1, document.BulletTypeDash)

	listn_3 := fmt.Sprintf("2.	%.1f%% имеют 3 и более взаимных выбора, из них %.1f%% положительных, - %.1f%% - отрицательных.",
		n3.n3_1, n3.n3_2, n3.n3_3)
	content = doc.AddNumberedList(listn_3, 0, document.ListTypeDecimal)

	listn_4 := fmt.Sprintf("%.1f%% обучающихся имеют противоречивые выборы. Ситуация противоречивых выборов болезненна и чревата негативными последствиями, особенно для  %.1f%% детей, которые адресуют положительный выбор.",
		n4.n4_1, n4.n4_2)
	content = doc.AddNumberedList(listn_4, 0, document.ListTypeDecimal)

	listn_5 := fmt.Sprintf("%.1f%% имеют адекватные представления о своем месте в группе, из них %.1f%% непринятых детей, что может вести к осознанной отгороженности от коллектива, замкнутости и не вовлеченности в общую деятельность. Часто такое явление имеет защитный характер, маскирует неуверенность, социальную тревожность ребенка.",
		n56.n5_1, n56.n5_2)
	content = doc.AddNumberedList(listn_5, 0, document.ListTypeDecimal)

	listn_6_1 := fmt.Sprintf("%.1f%% имеют неадекватные представления о своем месте в группе:", n56.n6_1)
	content = doc.AddNumberedList(listn_6_1, 0, document.ListTypeDecimal)

	n6_2 := fmt.Sprintf("%.1f%% непринятых детей из целевых групп, что может свидетельствовать о низком уровне умственного развития, инфантилизме, личностных нарушениях;",
		n56.n6_2)
	n6_3 := fmt.Sprintf("%.1f%% непринятых детей без личностных нарушений, что может свидетельствовать о защитной реакции, ухода от травмирующего воздействия и связаны с демонстрацией ребенком своей позиции или поведения, которое не соответствует ожиданиям группы;",
		n56.n6_3)
	n6_4 := fmt.Sprintf("%.1f%% принятых детей, что может говорить о внутренних конфликтах и эмоционально-личностных проблемах ребенка и проявляться одинаково болезненной реакции как на отсутствие ожидаемой от других агрессии в случае необоснованного субъективного занижения своего статуса, так и на отсутствие доброжелательности и поддержки при переоценке своей, роли в коллективе.",
		n56.n6_4)
	content = doc.AddBulletList(n6_2, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n6_3, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n6_4, 1, document.BulletTypeDash)

	listn_7 := ""
	if len(n7.n7) == 0 {
		listn_7 = "Референтная группа отсутствует."
		content = doc.AddNumberedList(listn_7, 0, document.ListTypeDecimal)
	} else {
		listn_7 = fmt.Sprintf("Референтная группа в сравнении с системой социометрических статусов:")
		content = doc.AddNumberedList(listn_7, 0, document.ListTypeDecimal)

		for _, val := range n7.n7 {
			content = doc.AddBulletList(val, 1, document.BulletTypeDash)
		}

	}
	content.SetStyle(style.StyleNormal)

	// 5. Save the document
	// fullPath := filepath.Join(SaveDir, "Аналитическая справка.docx")
	// fmt.Println(fullPath)
	err := doc.Save(fullPath)
	if err != nil {
		return fmt.Errorf("Failed to save document: %v", err)
	}

	return nil
}
