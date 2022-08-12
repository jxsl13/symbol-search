package main

import (
	"io"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
)

type SyncTable struct {
	mu sync.Mutex
	t  table.Writer
}

func (t *SyncTable) AppendFooter(row table.Row, configs ...table.RowConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.AppendFooter(row, configs...)
}
func (t *SyncTable) AppendHeader(row table.Row, configs ...table.RowConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.AppendHeader(row, configs...)
}
func (t *SyncTable) AppendRow(row table.Row, configs ...table.RowConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.AppendRow(row, configs...)
}
func (t *SyncTable) AppendRows(rows []table.Row, configs ...table.RowConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.AppendRows(rows, configs...)
}
func (t *SyncTable) AppendSeparator() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.AppendSeparator()
}
func (t *SyncTable) Length() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.t.Length()
}
func (t *SyncTable) Render() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.t.Render()
}
func (t *SyncTable) RenderCSV() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.t.RenderCSV()
}
func (t *SyncTable) RenderHTML() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.t.RenderHTML()
}
func (t *SyncTable) RenderMarkdown() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.t.RenderMarkdown()
}
func (t *SyncTable) ResetFooters() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.ResetFooters()
}
func (t *SyncTable) ResetHeaders() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.ResetHeaders()
}
func (t *SyncTable) ResetRows() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.ResetRows()
}
func (t *SyncTable) SetAllowedRowLength(length int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetAllowedRowLength(length)
}
func (t *SyncTable) SetAutoIndex(autoIndex bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetAutoIndex(autoIndex)
}
func (t *SyncTable) SetCaption(format string, a ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetCaption(format, a...)
}
func (t *SyncTable) SetColumnConfigs(configs []table.ColumnConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetColumnConfigs(configs)
}
func (t *SyncTable) SetIndexColumn(colNum int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetIndexColumn(colNum)
}
func (t *SyncTable) SetOutputMirror(mirror io.Writer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetOutputMirror(mirror)
}
func (t *SyncTable) SetPageSize(numLines int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetPageSize(numLines)
}
func (t *SyncTable) SetRowPainter(painter table.RowPainter) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetRowPainter(painter)
}
func (t *SyncTable) SetStyle(style table.Style) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetStyle(style)
}
func (t *SyncTable) SetTitle(format string, a ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SetTitle(format, a...)
}
func (t *SyncTable) SortBy(sortBy []table.SortBy) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SortBy(sortBy)
}
func (t *SyncTable) Style() *SyncStyle {
	t.mu.Lock()
	defer t.mu.Unlock()
	return &SyncStyle{
		mu: &t.mu,
		s:  t.t.Style(),
	}
}
func (t *SyncTable) SuppressEmptyColumns() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.t.SuppressEmptyColumns()
}

type SyncStyle struct {
	s  *table.Style
	mu *sync.Mutex
}

func (s *SyncStyle) Name() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Name
}
func (s *SyncStyle) Box() table.BoxStyle {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Box
}
func (s *SyncStyle) Color() table.ColorOptions {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Color
}
func (s *SyncStyle) Format() table.FormatOptions {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Format
}
func (s *SyncStyle) HTML() table.HTMLOptions {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.HTML
}
func (s *SyncStyle) Options() table.Options {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Options
}
func (s *SyncStyle) Title() table.TitleOptions {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Title
}

func (s *SyncStyle) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Name = name
}
func (s *SyncStyle) SetBox(box table.BoxStyle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Box = box
}
func (s *SyncStyle) SetColor(color table.ColorOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Color = color
}
func (s *SyncStyle) SetFormat(format table.FormatOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Format = format
}
func (s *SyncStyle) SetHTML(html table.HTMLOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.HTML = html
}
func (s *SyncStyle) SetOptions(options table.Options) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Options = options
}
func (s *SyncStyle) SetTitle(title table.TitleOptions) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Title = title
}
