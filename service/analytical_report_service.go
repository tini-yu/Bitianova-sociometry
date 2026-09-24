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
	subtitle := doc.AddParagraph("Аналитическая справка по результатам социометрии")
	subtitle.SetStyle(style.StyleSubtitle)

	// n1
	n1_1 := fmt.Sprintf("Дата проведения диагностики – %s",
		n1.n1_1)
	n1_2 := fmt.Sprintf("Кол-во участников диагностики – %d", n1.n1_2)
	content := doc.AddParagraph(n1_1)
	content = doc.AddParagraph(n1_2)

	content.SetStyle(style.StyleNormal)
	// n2
	content = doc.AddParagraph("Получены следующие эмпирические данные:\n")
	content.SetStyle(style.StyleNormal)

	listn_2_1 := fmt.Sprintf("В классе %.1f%% (%d чел.) принятых детей, из них %.1f%% (%d чел.) звезд, %.1f%% (%d чел.) - популярных. %.1f%% (%d чел.) непринятых детей, из них:\n",
		n2.n2_1_1, n2.n2_1_1_p, n2.n2_1_2, n2.n2_1_2_p, n2.n2_1_3, n2.n2_1_3_p, n2.n2_1_4, n2.n2_1_4_p)
	content = doc.AddNumberedList(listn_2_1, 0, document.ListTypeDecimal)
	content.SetStyle(style.StyleNormal)

	n2_2 := fmt.Sprintf("%.1f%% (%d чел.) отверженных детей, которых воспринимают весьма экспрессивно, но отрицательно, определенное неприятие человека, его качеств, свойств и привычек;", n2.n2_2, n2.n2_2_p)
	n2_3 := fmt.Sprintf("%.1f%% (%d чел.) изолированных детей, которых нет в эмоциональном реестре группы ни на уровне чувств, ни на уровне отношений;", n2.n2_3, n2.n2_3_p)
	n2_4 := fmt.Sprintf("%.1f%% (%d чел.) обучающихся непринятых детей из целевых групп, что может свидетельствовать о дефицитах в интеллектуальных, личностных, физических или материальных ресурсов и являться препятствием в расположении одноклассников;", n2.n2_4, n2.n2_4_p)
	content = doc.AddBulletList(n2_2, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n2_3, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n2_4, 1, document.BulletTypeDash)

	listn_3 := fmt.Sprintf("%.1f%% (%d чел.) имеют 3 и более взаимных выбора, из них %.1f%% (%d чел.) - положительных, %.1f%% (%d чел.) - отрицательных.",
		n3.n3_1, n3.n3_1_p, n3.n3_2, n3.n3_2_p, n3.n3_3, n3.n3_3_p)
	content = doc.AddNumberedList(listn_3, 0, document.ListTypeDecimal)

	listn_4 := fmt.Sprintf("%.1f%% (%d чел.) обучающихся имеют противоречивые выборы. Ситуация противоречивых выборов болезненна и чревата негативными последствиями, особенно для %.1f%% (%d чел.) детей, которые адресуют положительный выбор.",
		n4.n4_1, n4.n4_1_p, n4.n4_2, n4.n4_2_p)
	content = doc.AddNumberedList(listn_4, 0, document.ListTypeDecimal)

	listn_5 := fmt.Sprintf("%.1f%% (%d чел.) имеют адекватные представления о своем месте в группе, из них %.1f%% (%d чел.) непринятых детей, что может вести к осознанной отгороженности от коллектива, замкнутости и не вовлеченности в общую деятельность. Часто такое явление имеет защитный характер, маскирует неуверенность, социальную тревожность ребенка.",
		n56.n5_1, n56.n5_1_p, n56.n5_2, n56.n5_2_p)
	content = doc.AddNumberedList(listn_5, 0, document.ListTypeDecimal)

	listn_6_1 := fmt.Sprintf("%.1f%% (%d чел.) имеют неадекватные представления о своем месте в группе:", n56.n6_1, n56.n6_1_p)
	content = doc.AddNumberedList(listn_6_1, 0, document.ListTypeDecimal)

	n6_2 := fmt.Sprintf("%.1f%% (%d чел.) непринятых детей из целевых групп, что может свидетельствовать о низком уровне умственного развития, инфантилизме, личностных нарушениях;",
		n56.n6_2, n56.n6_2_p)
	n6_3 := fmt.Sprintf("%.1f%% (%d чел.) непринятых детей без личностных нарушений, что может свидетельствовать о защитной реакции, ухода от травмирующего воздействия и связаны с демонстрацией ребенком своей позиции или поведения, которое не соответствует ожиданиям группы;",
		n56.n6_3, n56.n6_3_p)
	n6_4 := fmt.Sprintf("%.1f%% (%d чел.) принятых детей, что может говорить о внутренних конфликтах и эмоционально-личностных проблемах ребенка и проявляться одинаково болезненной реакции как на отсутствие ожидаемой от других агрессии в случае необоснованного субъективного занижения своего статуса, так и на отсутствие доброжелательности и поддержки при переоценке своей, роли в коллективе.",
		n56.n6_4, n56.n6_4_p)
	content = doc.AddBulletList(n6_2, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n6_3, 1, document.BulletTypeDash)
	content = doc.AddBulletList(n6_4, 1, document.BulletTypeDash)

	listn_7 := ""
	if len(n7.n7) == 0 {
		listn_7 = "Референтная группа отсутствует."
		content = doc.AddNumberedList(listn_7, 0, document.ListTypeDecimal)
	} else {
		listn_7 = "Референтная группа в сравнении с системой социометрических статусов (рейтинг по референтности):"
		content = doc.AddNumberedList(listn_7, 0, document.ListTypeDecimal)

		for _, val := range n7.n7 {
			content = doc.AddBulletList(val, 1, document.BulletTypeDash)
		}

	}
	content.SetStyle(style.StyleNormal)

	// Выводы титульник:
	subtitleResults := doc.AddParagraph("На основе эмпирических данных можно сделать следующие выводы:")
	subtitleResults.SetStyle(style.StyleSubtitle)

	// n8
	n8_1 := "Преобладание или просто большое количество непринятых детей, говорит о низком уровне психологического развития группы или о ее распаде."
	n8_2 := "Преобладание принятых детей говорит о комфортной психологической атмосфере в группе."
	n8_3 := "Ситуация существования нескольких изгоев на фоне достаточно благоприятной эмоциональной структуры отношений других членов группы свидетельствует о типе «сплоченный класс против одного».  Решение конфликтов в этой группе происходит за счет «козла отпущения». Даная ситуация так или иначе связано с «руководящей линией» взрослого, работающего с классом."
	n8_4 := "В данной группе всем достается некоторое количество отвержений. Решение конфликтов в коллективе происходит через поиск причин и решений."
	n8_5 := "Наличие внутри коллектива большого количества противоречивых связей свидетельствует о разрозненной групповой структуре, что говорит о неблагополучии эмоциональных отношений, возможном распаде или близком переструктурировании отношений. Группа пребывает на сложной ступени своего развития, характеризующейся и конфликтностью, и низкой групповой эмпатией, и эмоциональной неудовлетворенностью."
	n8_6 := "Эмоциональная привлекательность значительной части референтной группы не совпадает с ценностной значимостью. В число референтометрических лидеров попадают изгои или отвергаемые, значит существует лабильность в группе."
	n8_7_1 := "Эмоциональная привлекательность значительной части референтной группы совпадает с ценностной значимостью, что говорит о том, что в классе существует открытость собственного мнения. Потенциальные лидеры:"

	n8_7_2 := "\tобладают реальным влиянием на мысли и поступки членов группы и в критических ситуациях за ними пойдут как за лидером."
	n8_8 := fmt.Sprintf("Рассчитывать на высокую учебную активность у %.1f%% (%d чел.) детей нецелесообразно, так как они глубоко погружены в свои эмоциональные отношения. При этом %.1f%% (%d чел.) детей с адекватной рефлексией и %.1f%% (%d чел.) с неадекватной рефлексией.",
		n8.n8_8_1, n8.n8_8_1_p, n8.n8_8_2, n8.n8_8_2_p, n8.n8_8_3, n8.n8_8_3_p)
	if n8.n8_1 {
		content = doc.AddBulletList(n8_1, 0, document.BulletTypeDash)
	}
	if n8.n8_2 {
		content = doc.AddBulletList(n8_2, 0, document.BulletTypeDash)
	}
	if n8.n8_3 {
		content = doc.AddBulletList(n8_3, 0, document.BulletTypeDash)
	}
	if n8.n8_4 {
		content = doc.AddBulletList(n8_4, 0, document.BulletTypeDash)
	}
	if n8.n8_5 {
		content = doc.AddBulletList(n8_5, 0, document.BulletTypeDash)
	}
	if n8.n8_6 {
		content = doc.AddBulletList(n8_6, 0, document.BulletTypeDash)
	}
	if n8.n8_7 {
		content = doc.AddBulletList(n8_7_1, 0, document.BulletTypeDash)
		for _, name := range n8.n8_7_list {
			content = doc.AddBulletList(name, 1, document.BulletTypeDot)
		}
		content = doc.AddBulletList(n8_7_2, 0, document.BulletTypeDash)
	}
	if n8.n8_8 {
		content = doc.AddBulletList(n8_8, 0, document.BulletTypeDash)
	}
	content.SetStyle(style.StyleNormal)

	// 5. Save the document
	err := doc.Save(fullPath)
	if err != nil {
		return fmt.Errorf("Failed to save document: %v", err)
	}

	return nil
}
