package handlers

import "testing"

func TestParseTraceCodeWordsByColumns(t *testing.T) {
	words := []baiduOCRWord{
		word("交易流水号", 544, 121),
		word("药品编码", 547, 321),
		word("药品名称", 549, 526),
		word("结算日期", 555, 934),
		word("就诊id", 561, 1460),
		word("011100050Y260518010818", 564, 64),
		word("ZF03AAG0012010401007", 566, 255),
		word("甘桔冰梅片", 568, 450),
		word("1.0000", 571, 668),
		word("1.0000", 573, 777),
		word("2026-05-1815:35:56", 574, 883),
		word("20260518153551098373", 581, 1387),
		word("0111000520001127410Y", 583, 64),
		word("XL02BGA045A001010204641", 585, 255),
		word("阿那曲唑片", 588, 449),
		word("2.0000", 592, 776),
		word("2026-05-2713:34:13", 594, 883),
		word("20260527133330579210", 600, 1386),
		word("0111000520001129583Y", 603, 63),
		word("ZA09BAD0254010102017", 605, 255),
		word("地榆升白片", 608, 449),
		word("4.0000", 612, 776),
		word("2026-06-0114:39:14", 613, 882),
		word("20260601143832593173", 620, 1385),
	}

	records := parseTraceCodeWords(words)
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d: %#v", len(records), records)
	}

	if records[0].TransactionSerialNumber != "011100050Y260518010818" {
		t.Fatalf("unexpected first serial: %s", records[0].TransactionSerialNumber)
	}
	if records[0].TransactionSerialNumber == "20260518153551098373" {
		t.Fatal("visit id was incorrectly parsed as transaction serial number")
	}
	if records[1].DrugCode != "XL02BGA045A001010204641" || records[1].DrugName != "阿那曲唑片" {
		t.Fatalf("unexpected second record: %#v", records[1])
	}
	if records[2].SettlementDate != "2026-06-01" {
		t.Fatalf("unexpected third settlement date: %s", records[2].SettlementDate)
	}
}

func word(text string, top int, left int) baiduOCRWord {
	return baiduOCRWord{
		Words: text,
		Location: baiduOCRLocation{
			Top:  top,
			Left: left,
		},
	}
}
