package report

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/secscan/secscan/backend/internal/domain"
)

func GeneratePDF(scan domain.Scan) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("SecScan Security Report", false)
	pdf.SetMargins(16, 16, 16)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "SecScan Security Report")
	pdf.Ln(14)

	pdf.SetFont("Arial", "", 11)
	writeLine(pdf, "Target URL", scan.CanonicalURL)
	writeLine(pdf, "Scan ID", scan.ID)
	writeLine(pdf, "Status", string(scan.Status))
	writeLine(pdf, "Overall Grade", fmt.Sprintf("%s (%d/100)", scan.OverallGrade, scan.OverallScore))
	writeLine(pdf, "Created", scan.CreatedAt.Format("2006-01-02 15:04:05 MST"))
	if scan.CompletedAt != nil {
		writeLine(pdf, "Completed", scan.CompletedAt.Format("2006-01-02 15:04:05 MST"))
	}

	pdf.Ln(6)
	section(pdf, "Module Summary")
	for _, module := range orderedModuleResults(scan) {
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(0, 7, fmt.Sprintf("%s: %s (%d/100)", strings.ToUpper(module.Name), module.Grade, module.Score))
		pdf.Ln(7)
		if module.Error != "" {
			pdf.SetFont("Arial", "", 10)
			pdf.MultiCell(0, 5, "Error: "+module.Error, "", "", false)
		}
	}

	pdf.Ln(4)
	section(pdf, "Key Findings")
	if len(scan.Findings) == 0 {
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, "No findings were produced by the selected scanner modules.", "", "", false)
	} else {
		limit := len(scan.Findings)
		if limit > 20 {
			limit = 20
		}
		for i := 0; i < limit; i++ {
			finding := scan.Findings[i]
			pdf.SetFont("Arial", "B", 10)
			pdf.MultiCell(0, 5, fmt.Sprintf("[%s] %s", strings.ToUpper(finding.Severity), finding.Title), "", "", false)
			pdf.SetFont("Arial", "", 9)
			pdf.MultiCell(0, 5, finding.Description, "", "", false)
			if finding.Evidence != "" {
				pdf.MultiCell(0, 5, "Evidence: "+finding.Evidence, "", "", false)
			}
			pdf.MultiCell(0, 5, "Recommendation: "+finding.Recommendation, "", "", false)
			pdf.Ln(2)
		}
	}

	pdf.Ln(4)
	section(pdf, "Recommendations")
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, "Treat this report as a lightweight academic assessment. Confirm any high-risk finding manually, keep public services patched, prefer HTTPS, and only scan systems you own or have explicit permission to test.", "", "", false)

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func writeLine(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(38, 6, label+":", "", 0, "", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, value, "", "", false)
}

func section(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, title)
	pdf.Ln(10)
}

func orderedModuleResults(scan domain.Scan) []domain.ModuleResult {
	results := make([]domain.ModuleResult, 0, len(scan.Results))
	for _, result := range scan.Results {
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	return results
}
