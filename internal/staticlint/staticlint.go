// Package staticlint реализует линтер и включает в себя следующие проверки:
//
// 1. базовые статические анализаторы пакета golang.org/x/tools/go/analysis/passes:
//   - printf
//   - errorsas
//   - nilness
//   - shadow
//   - unusedresult
//   - copylock
//   - composite
//   - httpresponse
//   - lostcancel
//   - loopclosure
//   - nilfunc
//
// 2. анализаторы класса SA пакета staticlint.io
//
// 3. анализаторы simple и quickfix пакета staticlint.io
//
// 4. публичные анализаторы github.com/timakin/bodyclose и github.com/jingyugao/rowserrcheck
//
// 5. кастомный анализатор osexitusage.
package staticlint

import (
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"github.com/mkolibaba/metrics/internal/staticlint/osexitusage"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func Run() {
	checks := []*analysis.Analyzer{
		// базовые статические анализаторы пакета golang.org/x/tools/go/analysis/passes
		printf.Analyzer,
		errorsas.Analyzer,
		nilness.Analyzer,
		shadow.Analyzer,
		unusedresult.Analyzer,
		copylock.Analyzer,
		composite.Analyzer,
		httpresponse.Analyzer,
		lostcancel.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,

		// публичные анализаторы
		bodyclose.Analyzer,
		rowserr.NewAnalyzer(
			"github.com/jackc/pgx/v5/stdlib",
		),

		// кастомный анализатор osexitusage
		osexitusage.Analyzer,
	}

	// анализаторы класса SA пакета staticlint.io
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	// анализаторы остальных классов пакета staticlint.io
	for _, v := range simple.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range quickfix.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	multichecker.Main(checks...)
}
