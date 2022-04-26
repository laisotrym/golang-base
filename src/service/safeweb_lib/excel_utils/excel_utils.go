package safeweb_lib_excel_utils

import (
    "errors"
    "reflect"
    "strconv"
    
    "github.com/360EntSecGroup-Skylar/excelize"
    
    "safeweb.app/service/campaign_merchant/types"
)

const (
    MinimalColumnCampaignMerchant = 7
)

func ReadExcelFileToSlice(filePath string, structType reflect.Type) (interface{}, error) {
    excelFile, err := excelize.OpenFile(filePath)
    if err != nil {
        return nil, err
    }
    sheetName := excelFile.GetSheetName(excelFile.GetActiveSheetIndex())
    
    rows := excelFile.GetRows(sheetName)
    for _, row := range rows {
        for _, colCell := range row {
            print(colCell, "\t")
        }
        println()
    }
    
    rs := reflect.MakeSlice(reflect.SliceOf(structType), len(rows), len(rows))
    return rs, nil
}

func ReadExcelFileToSliceOfCampaignMerchant(filePath string) ([]types.CampaignMerchantDTO, error) {
    excelFile, err := excelize.OpenFile(filePath)
    if err != nil {
        return nil, err
    }
    sheetName := excelFile.GetSheetName(excelFile.GetActiveSheetIndex())
    
    rows := excelFile.GetRows(sheetName)
    rs := make([]types.CampaignMerchantDTO, 0, len(rows))
    
    for _, row := range rows {
        temp, err := readOneRowToCampaignMerchant(row)
        if err == nil {
            rs = append(rs, *temp)
        }
    }
    return rs, nil
}

func readOneRowToCampaignMerchant(row []string) (rs *types.CampaignMerchantDTO, err error) {
    if len(row) < MinimalColumnCampaignMerchant {
        return nil, errors.New("Invalid row.")
    }
    no, err := strconv.Atoi(row[0])
    if err != nil {
        return nil, err
    }
    return &types.CampaignMerchantDTO{
        No:           no,
        MerchantCode: row[1],
        TerminalCode: row[3],
        StartTime:    row[5],
        EndTime:      row[6],
    }, nil
}
