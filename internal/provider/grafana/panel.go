package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

// Each panel may be one of these types.
const (
	CustomType panelType = iota
	DashlistType
	GraphType
	TableType
	TextType
	PluginlistType
	AlertlistType
	SinglestatType
	StatType
	RowType
	BarGaugeType
	GaugeType
	HeatmapType
	TimeseriesType
)

type (
	// Panel represents panels of different types defined in Grafana.
	Panel struct {
		CommonPanel
		// Should be initialized only one type of panels.
		// OfType field defines which of types below will be used.
		*GraphPanel
		*TablePanel
		*TextPanel
		*SinglestatPanel
		*StatPanel
		*DashlistPanel
		*PluginlistPanel
		*RowPanel
		*AlertlistPanel
		*BarGaugePanel
		*GaugePanel
		*HeatmapPanel
		*TimeseriesPanel
		*CustomPanel
	}
	panelType   int8
	CommonPanel struct {
		Datasource interface{} `json:"datasource,omitempty"` // metrics
		Editable   bool        `json:"editable"`
		Error      bool        `json:"error"`
		GridPos    struct {
			H *int `json:"h,omitempty"`
			W *int `json:"w,omitempty"`
			X *int `json:"x,omitempty"`
			Y *int `json:"y,omitempty"`
		} `json:"gridPos,omitempty"`
		Height           interface{} `json:"height,omitempty"` // general
		HideTimeOverride *bool       `json:"hideTimeOverride,omitempty"`
		ID               uint        `json:"id"`
		IsNew            bool        `json:"isNew"`
		//Links            []Link      `json:"links,omitempty"`    // general
		MinSpan  *float32  `json:"minSpan,omitempty"`  // templating options
		OfType   panelType `json:"-"`                  // it required for defining type of the panel
		Renderer *string   `json:"renderer,omitempty"` // display styles
		Repeat   *string   `json:"repeat,omitempty"`   // templating options
		// RepeatIteration *int64   `json:"repeatIteration,omitempty"`
		RepeatPanelID *uint `json:"repeatPanelId,omitempty"`
		ScopedVars    map[string]struct {
			Selected bool   `json:"selected"`
			Text     string `json:"text"`
			Value    string `json:"value"`
		} `json:"scopedVars,omitempty"`
		Span        float32 `json:"span"`                  // general
		Title       string  `json:"title"`                 // general
		Description *string `json:"description,omitempty"` // general
		Transparent bool    `json:"transparent"`
		Type        string  `json:"type"`
		Alert       *Alert  `json:"alert,omitempty"`
	}
	AlertEvaluator struct {
		Params []float64 `json:"params,omitempty"`
		Type   string    `json:"type,omitempty"`
	}
	AlertOperator struct {
		Type string `json:"type,omitempty"`
	}
	AlertQuery struct {
		Params []string `json:"params,omitempty"`
	}
	AlertReducer struct {
		Params []string `json:"params,omitempty"`
		Type   string   `json:"type,omitempty"`
	}
	AlertCondition struct {
		Evaluator AlertEvaluator `json:"evaluator,omitempty"`
		Operator  AlertOperator  `json:"operator,omitempty"`
		Query     AlertQuery     `json:"query,omitempty"`
		Reducer   AlertReducer   `json:"reducer,omitempty"`
		Type      string         `json:"type,omitempty"`
	}
	Alert struct {
		AlertRuleTags       map[string]string `json:"alertRuleTags,omitempty"`
		Conditions          []AlertCondition  `json:"conditions,omitempty"`
		ExecutionErrorState string            `json:"executionErrorState,omitempty"`
		Frequency           string            `json:"frequency,omitempty"`
		Handler             int               `json:"handler,omitempty"`
		Name                string            `json:"name,omitempty"`
		NoDataState         string            `json:"noDataState,omitempty"`
		//Notifications       []AlertNotification `json:"notifications,omitempty"`
		Message string `json:"message,omitempty"`
		For     string `json:"for,omitempty"`
	}
	GraphPanel struct {
		AliasColors interface{} `json:"aliasColors"` // XXX
		Bars        bool        `json:"bars"`
		DashLength  *uint       `json:"dashLength,omitempty"`
		Dashes      *bool       `json:"dashes,omitempty"`
		Decimals    *int        `json:"decimals,omitempty"`
		Fill        int         `json:"fill"`
		//		Grid        grid        `json:"grid"` obsoleted in 4.1 by xaxis and yaxis

		Legend          Legend           `json:"legend,omitempty"`
		LeftYAxisLabel  *string          `json:"leftYAxisLabel,omitempty"`
		Lines           bool             `json:"lines"`
		Linewidth       uint             `json:"linewidth"`
		NullPointMode   string           `json:"nullPointMode"`
		Percentage      bool             `json:"percentage"`
		Pointradius     float32          `json:"pointradius"`
		Points          bool             `json:"points"`
		RightYAxisLabel *string          `json:"rightYAxisLabel,omitempty"`
		SeriesOverrides []SeriesOverride `json:"seriesOverrides,omitempty"`
		SpaceLength     *uint            `json:"spaceLength,omitempty"`
		Stack           bool             `json:"stack"`
		SteppedLine     bool             `json:"steppedLine"`
		Targets         []Target         `json:"targets,omitempty"`
		Thresholds      []Threshold      `json:"thresholds,omitempty"`
		TimeFrom        *string          `json:"timeFrom,omitempty"`
		TimeShift       *string          `json:"timeShift,omitempty"`
		Tooltip         Tooltip          `json:"tooltip"`
		XAxis           bool             `json:"x-axis,omitempty"`
		YAxis           bool             `json:"y-axis,omitempty"`
		YFormats        []string         `json:"y_formats,omitempty"`
		Xaxis           Axis             `json:"xaxis"` // was added in Grafana 4.x?
		Yaxes           []Axis           `json:"yaxes"` // was added in Grafana 4.x?
		FieldConfig     *FieldConfig     `json:"fieldConfig,omitempty"`
	}
	FieldConfig struct {
		Defaults FieldConfigDefaults `json:"defaults"`
	}
	TextSize struct {
		TitleSize *int `json:"titleSize,omitempty"`
		ValueSize *int `json:"valueSize,omitempty"`
	}
	ReduceOptions struct {
		Values bool     `json:"values"`
		Fields string   `json:"fields"`
		Limit  *int     `json:"limit,omitempty"`
		Calcs  []string `json:"calcs"`
	}
	Options struct {
		Orientation string `json:"orientation"`
		TextMode    string `json:"textMode"`
		ColorMode   string `json:"colorMode"`
		GraphMode   string `json:"graphMode"`
		JustifyMode string `json:"justifyMode"`
		DisplayMode string `json:"displayMode"`
		Content     string `json:"content"`
		Mode        string `json:"mode"`
		// gauge specific
		ShowThresholdLabels  *bool `json:"showThresholdLabels,omitempty"`
		ShowThresholdMarkers *bool `json:"showThresholdMarkers,omitempty"`
		// etc
		TextSize      TextSize      `json:"text"`
		ReduceOptions ReduceOptions `json:"reduceOptions"`
	}
	Threshold struct {
		// the alert threshold value, we do not omitempty, since 0 is a valid
		// threshold
		Value float32 `json:"value"`
		// critical, warning, ok, custom
		ColorMode string `json:"colorMode,omitempty"`
		// gt or lt
		Op   string `json:"op,omitempty"`
		Fill bool   `json:"fill"`
		Line bool   `json:"line"`
		// hexadecimal color (e.g. #629e51, only when ColorMode is "custom")
		FillColor string `json:"fillColor,omitempty"`
		// hexadecimal color (e.g. #629e51, only when ColorMode is "custom")
		LineColor string `json:"lineColor,omitempty"`
		// left or right
		Yaxis string `json:"yaxis,omitempty"`
	}

	Tooltip struct {
		Shared       bool   `json:"shared"`
		ValueType    string `json:"value_type"`
		MsResolution bool   `json:"msResolution,omitempty"` // was added in Grafana 3.x
		Sort         int    `json:"sort,omitempty"`
	}
	TablePanel struct {
		Columns   []Column      `json:"columns"`
		Sort      *Sort         `json:"sort,omitempty"`
		Styles    []ColumnStyle `json:"styles"`
		Transform string        `json:"transform"`
		Targets   []Target      `json:"targets,omitempty"`
		Scroll    bool          `json:"scroll"` // from grafana 3.x
	}
	TextPanel struct {
		Content     string        `json:"content"`
		Mode        string        `json:"mode"`
		PageSize    uint          `json:"pageSize"`
		Scroll      bool          `json:"scroll"`
		ShowHeader  bool          `json:"showHeader"`
		Sort        Sort          `json:"sort"`
		Styles      []ColumnStyle `json:"styles"`
		FieldConfig FieldConfig   `json:"fieldConfig"`
		Options     struct {
			Content string `json:"content"`
			Mode    string `json:"mode"`
		} `json:"options"`
	}
	SinglestatPanel struct {
		Colors          []string   `json:"colors"`
		ColorValue      bool       `json:"colorValue"`
		ColorBackground bool       `json:"colorBackground"`
		Decimals        int        `json:"decimals"`
		Format          string     `json:"format"`
		Gauge           Gauge      `json:"gauge,omitempty"`
		MappingType     *uint      `json:"mappingType,omitempty"`
		MappingTypes    []*MapType `json:"mappingTypes,omitempty"`
		//MaxDataPoints   *IntString  `json:"maxDataPoints,omitempty"`
		NullPointMode   string      `json:"nullPointMode"`
		Postfix         *string     `json:"postfix,omitempty"`
		PostfixFontSize *string     `json:"postfixFontSize,omitempty"`
		Prefix          *string     `json:"prefix,omitempty"`
		PrefixFontSize  *string     `json:"prefixFontSize,omitempty"`
		RangeMaps       []*RangeMap `json:"rangeMaps,omitempty"`
		SparkLine       SparkLine   `json:"sparkline,omitempty"`
		Targets         []Target    `json:"targets,omitempty"`
		Thresholds      string      `json:"thresholds"`
		ValueFontSize   string      `json:"valueFontSize"`
		ValueMaps       []ValueMap  `json:"valueMaps"`
		ValueName       string      `json:"valueName"`
	}
	GaugePanel struct {
		Options     Options     `json:"options"`
		Targets     []Target    `json:"targets,omitempty"`
		FieldConfig FieldConfig `json:"fieldConfig"`
	}
	StatPanel struct {
		Colors          []string   `json:"colors"`
		ColorValue      bool       `json:"colorValue"`
		ColorBackground bool       `json:"colorBackground"`
		Decimals        int        `json:"decimals"`
		Format          string     `json:"format"`
		Gauge           Gauge      `json:"gauge,omitempty"`
		MappingType     *uint      `json:"mappingType,omitempty"`
		MappingTypes    []*MapType `json:"mappingTypes,omitempty"`
		//MaxDataPoints   *IntString  `json:"maxDataPoints,omitempty"`
		NullPointMode   string      `json:"nullPointMode"`
		Postfix         *string     `json:"postfix,omitempty"`
		PostfixFontSize *string     `json:"postfixFontSize,omitempty"`
		Prefix          *string     `json:"prefix,omitempty"`
		PrefixFontSize  *string     `json:"prefixFontSize,omitempty"`
		RangeMaps       []*RangeMap `json:"rangeMaps,omitempty"`
		SparkLine       SparkLine   `json:"sparkline,omitempty"`
		Targets         []Target    `json:"targets,omitempty"`
		Thresholds      string      `json:"thresholds"`
		ValueFontSize   string      `json:"valueFontSize"`
		ValueMaps       []ValueMap  `json:"valueMaps"`
		ValueName       string      `json:"valueName"`
		Options         Options     `json:"options"`
		FieldConfig     FieldConfig `json:"fieldConfig"`
	}
	DashlistPanel struct {
		Mode     string   `json:"mode"`
		Query    string   `json:"query"`
		Tags     []string `json:"tags"`
		FolderID int      `json:"folderId"`
		Limit    int      `json:"limit"`
		Headings bool     `json:"headings"`
		Recent   bool     `json:"recent"`
		Search   bool     `json:"search"`
		Starred  bool     `json:"starred"`
	}
	PluginlistPanel struct {
		Limit int `json:"limit,omitempty"`
	}
	AlertlistPanel struct {
		OnlyAlertsOnDashboard bool     `json:"onlyAlertsOnDashboard"`
		Show                  string   `json:"show"`
		SortOrder             int      `json:"sortOrder"`
		Limit                 int      `json:"limit"`
		StateFilter           []string `json:"stateFilter"`
		NameFilter            string   `json:"nameFilter,omitempty"`
		DashboardTags         []string `json:"dashboardTags,omitempty"`
	}
	BarGaugePanel struct {
		Options     Options     `json:"options"`
		Targets     []Target    `json:"targets,omitempty"`
		FieldConfig FieldConfig `json:"fieldConfig"`
	}
	RowPanel struct {
		Panels    []Panel `json:"panels"`
		Collapsed bool    `json:"collapsed"`
	}
	HeatmapPanel struct {
		Cards struct {
			CardPadding *float64 `json:"cardPadding"`
			CardRound   *float64 `json:"cardRound"`
		} `json:"cards"`
		Color struct {
			CardColor   string   `json:"cardColor"`
			ColorScale  string   `json:"colorScale"`
			ColorScheme string   `json:"colorScheme"`
			Exponent    float64  `json:"exponent"`
			Min         *float64 `json:"min,omitempty"`
			Max         *float64 `json:"max,omitempty"`
			Mode        string   `json:"mode"`
		} `json:"color"`
		DataFormat      string `json:"dataFormat"`
		HideZeroBuckets bool   `json:"hideZeroBuckets"`
		HighlightCards  bool   `json:"highlightCards"`
		Legend          struct {
			Show bool `json:"show"`
		} `json:"legend"`
		ReverseYBuckets bool     `json:"reverseYBuckets"`
		Targets         []Target `json:"targets,omitempty"`
		Tooltip         struct {
			Show          bool `json:"show"`
			ShowHistogram bool `json:"showHistogram"`
		} `json:"tooltip"`
		TooltipDecimals int `json:"tooltipDecimals"`
		XAxis           struct {
			Show bool `json:"show"`
		} `json:"xAxis"`
		XBucketNumber *float64 `json:"xBucketNumber"`
		XBucketSize   *string  `json:"xBucketSize"`
		YAxis         struct {
			Decimals    *int     `json:"decimals"`
			Format      string   `json:"format"`
			LogBase     int      `json:"logBase"`
			Show        bool     `json:"show"`
			Max         *string  `json:"max"`
			Min         *string  `json:"min"`
			SplitFactor *float64 `json:"splitFactor"`
		} `json:"yAxis"`
		YBucketBound  string   `json:"yBucketBound"`
		YBucketNumber *float64 `json:"yBucketNumber"`
		YBucketSize   *float64 `json:"yBucketSize"`
	}
	TimeseriesPanel struct {
		Targets     []Target          `json:"targets,omitempty"`
		Options     TimeseriesOptions `json:"options"`
		FieldConfig FieldConfig       `json:"fieldConfig"`
	}
	TimeseriesOptions struct {
		Legend  TimeseriesLegendOptions  `json:"legend,omitempty"`
		Tooltip TimeseriesTooltipOptions `json:"tooltip,omitempty"`
	}
	TimeseriesLegendOptions struct {
		Calcs       []string `json:"calcs"`
		DisplayMode string   `json:"displayMode"`
		Placement   string   `json:"placement"`
	}
	TimeseriesTooltipOptions struct {
		Mode string `json:"mode"`
	}
	FieldConfigDefaults struct {
		Unit       string            `json:"unit"`
		Decimals   *int              `json:"decimals,omitempty"`
		Min        *float64          `json:"min,omitempty"`
		Max        *float64          `json:"max,omitempty"`
		NoValue    *float64          `json:"noValue,omitempty"`
		Color      FieldConfigColor  `json:"color"`
		Thresholds Thresholds        `json:"thresholds"`
		Custom     FieldConfigCustom `json:"custom"`
		Mappings   []FieldMapping    `json:"mappings,omitempty"`
		//Links      []Link            `json:"links,omitempty"`
	}
	FieldMapping struct {
		Type    string                 `json:"type"`
		Options map[string]interface{} `json:"options"`
	}
	FieldConfigCustom struct {
		AxisLabel         string `json:"axisLabel,omitempty"`
		AxisPlacement     string `json:"axisPlacement"`
		AxisSoftMin       *int   `json:"axisSoftMin,omitempty"`
		AxisSoftMax       *int   `json:"axisSoftMax,omitempty"`
		BarAlignment      int    `json:"barAlignment"`
		DrawStyle         string `json:"drawStyle"`
		FillOpacity       int    `json:"fillOpacity"`
		GradientMode      string `json:"gradientMode"`
		LineInterpolation string `json:"lineInterpolation"`
		LineWidth         int    `json:"lineWidth"`
		PointSize         int    `json:"pointSize"`
		ShowPoints        string `json:"showPoints"`
		SpanNulls         bool   `json:"spanNulls"`
		HideFrom          struct {
			Legend  bool `json:"legend"`
			Tooltip bool `json:"tooltip"`
			Viz     bool `json:"viz"`
		} `json:"hideFrom"`
		LineStyle struct {
			Fill string `json:"fill"`
		} `json:"lineStyle"`
		ScaleDistribution struct {
			Type string `json:"type"`
			Log  int    `json:"log,omitempty"`
		} `json:"scaleDistribution"`
		Stacking struct {
			Group string `json:"group"`
			Mode  string `json:"mode"`
		} `json:"stacking"`
		ThresholdsStyle struct {
			Mode string `json:"mode"`
		} `json:"thresholdsStyle"`
	}
	Thresholds struct {
		Mode  string          `json:"mode"`
		Steps []ThresholdStep `json:"steps"`
	}
	ThresholdStep struct {
		Color string   `json:"color"`
		Value *float64 `json:"value"`
	}
	FieldConfigColor struct {
		Mode       string `json:"mode"`
		FixedColor string `json:"fixedColor,omitempty"`
		SeriesBy   string `json:"seriesBy,omitempty"`
	}
	CustomPanel map[string]interface{}
)

// for a graph panel
type (
	Axis struct {
		Format   string   `json:"format"`
		LogBase  int      `json:"logBase"`
		Decimals int      `json:"decimals,omitempty"`
		Max      *float64 `json:"max,omitempty"`
		Min      *float64 `json:"min,omitempty"`
		Show     bool     `json:"show"`
		Label    string   `json:"label,omitempty"`
	}
	SeriesOverride struct {
		Alias         string  `json:"alias"`
		Bars          *bool   `json:"bars,omitempty"`
		Color         *string `json:"color,omitempty"`
		Dashes        *bool   `json:"dashes,omitempty"`
		Fill          *int    `json:"fill,omitempty"`
		FillBelowTo   *string `json:"fillBelowTo,omitempty"`
		Legend        *bool   `json:"legend,omitempty"`
		Lines         *bool   `json:"lines,omitempty"`
		LineWidth     *int    `json:"linewidth,omitempty"`
		Stack         *bool   `json:"stack,omitempty"`
		Transform     *string `json:"transform,omitempty"`
		YAxis         *int    `json:"yaxis,omitempty"`
		ZIndex        *int    `json:"zindex,omitempty"`
		NullPointMode *string `json:"nullPointMode,omitempty"`
	}
	Sort struct {
		Col  int  `json:"col"`
		Desc bool `json:"desc"`
	}
	Legend struct {
		AlignAsTable bool  `json:"alignAsTable"`
		Avg          bool  `json:"avg"`
		Current      bool  `json:"current"`
		HideEmpty    bool  `json:"hideEmpty"`
		HideZero     bool  `json:"hideZero"`
		Max          bool  `json:"max"`
		Min          bool  `json:"min"`
		RightSide    bool  `json:"rightSide"`
		Show         bool  `json:"show"`
		SideWidth    *uint `json:"sideWidth,omitempty"`
		Total        bool  `json:"total"`
		Values       bool  `json:"values"`
	}
)

// for a table
type (
	Column struct {
		TextType string `json:"text"`
		Value    string `json:"value"`
	}
	ColumnStyle struct {
		Alias           *string    `json:"alias"`
		DateFormat      *string    `json:"dateFormat,omitempty"`
		Pattern         string     `json:"pattern"`
		Type            string     `json:"type"`
		ColorMode       *string    `json:"colorMode,omitempty"`
		Colors          *[]string  `json:"colors,omitempty"`
		Decimals        *int       `json:"decimals,omitempty"`
		Thresholds      *[]string  `json:"thresholds,omitempty"`
		Unit            *string    `json:"unit,omitempty"`
		MappingType     int        `json:"mappingType,omitempty"`
		ValueMaps       []ValueMap `json:"valueMaps,omitempty"`
		Link            bool       `json:"link,omitempty"`
		LinkTooltip     *string    `json:"linkTooltip,omitempty"`
		LinkUrl         *string    `json:"linkUrl,omitempty"`
		LinkTargetBlank bool       `json:"linkTargetBlank,omitempty"`
	}
)

// for a stat
type (
	ValueMap struct {
		Op       string `json:"op"`
		TextType string `json:"text"`
		Value    string `json:"value"`
	}
	Gauge struct {
		MaxValue         float32 `json:"maxValue"`
		MinValue         float32 `json:"minValue"`
		Show             bool    `json:"show"`
		ThresholdLabels  bool    `json:"thresholdLabels"`
		ThresholdMarkers bool    `json:"thresholdMarkers"`
	}
	SparkLine struct {
		FillColor *string  `json:"fillColor,omitempty"`
		Full      bool     `json:"full,omitempty"`
		LineColor *string  `json:"lineColor,omitempty"`
		Show      bool     `json:"show,omitempty"`
		YMin      *float64 `json:"ymin,omitempty"`
		YMax      *float64 `json:"ymax,omitempty"`
	}
)

// for an any panel
type Target struct {
	RefID      string      `json:"refId"`
	Datasource interface{} `json:"datasource,omitempty"`
	Hide       bool        `json:"hide,omitempty"`

	// For Prometheus
	Expr           string `json:"expr,omitempty"`
	IntervalFactor int    `json:"intervalFactor,omitempty"`
	Interval       string `json:"interval,omitempty"`
	Step           int    `json:"step,omitempty"`
	LegendFormat   string `json:"legendFormat,omitempty"`
	Instant        bool   `json:"instant,omitempty"`
	Format         string `json:"format,omitempty"`

	// For Graphite
	Target string `json:"target,omitempty"`

	// For CloudWatch
	Namespace  string            `json:"namespace,omitempty"`
	MetricName string            `json:"metricName,omitempty"`
	Statistics []string          `json:"statistics,omitempty"`
	Dimensions map[string]string `json:"dimensions,omitempty"`
	Period     string            `json:"period,omitempty"`
	Region     string            `json:"region,omitempty"`
	Label      string            `json:"label,omitempty"`
}

type MapType struct {
	Name  *string `json:"name,omitempty"`
	Value *int    `json:"value,omitempty"`
}

type RangeMap struct {
	From *string `json:"from,omitempty"`
	Text *string `json:"text,omitempty"`
	To   *string `json:"to,omitempty"`
}

type probePanel struct {
	CommonPanel
	//	json.RawMessage
}

func (p *Panel) UnmarshalJSON(b []byte) (err error) {
	var probe probePanel
	if err = json.Unmarshal(b, &probe); err != nil {
		return err
	}

	p.CommonPanel = probe.CommonPanel
	switch probe.Type {
	case "graph":
		var graph GraphPanel
		p.OfType = GraphType
		if err = json.Unmarshal(b, &graph); err == nil {
			p.GraphPanel = &graph
		}
	case "table":
		var table TablePanel
		p.OfType = TableType
		if err = json.Unmarshal(b, &table); err == nil {
			p.TablePanel = &table
		}
	case "text":
		var text TextPanel
		p.OfType = TextType
		if err = json.Unmarshal(b, &text); err == nil {
			p.TextPanel = &text
		}
	case "singlestat":
		var singlestat SinglestatPanel
		p.OfType = SinglestatType
		if err = json.Unmarshal(b, &singlestat); err == nil {
			p.SinglestatPanel = &singlestat
		}
	case "stat":
		var stat StatPanel
		p.OfType = StatType
		if err = json.Unmarshal(b, &stat); err == nil {
			p.StatPanel = &stat
		}
	case "dashlist":
		var dashlist DashlistPanel
		p.OfType = DashlistType
		if err = json.Unmarshal(b, &dashlist); err == nil {
			p.DashlistPanel = &dashlist
		}
	case "bargauge":
		var bargauge BarGaugePanel
		p.OfType = BarGaugeType
		if err = json.Unmarshal(b, &bargauge); err == nil {
			p.BarGaugePanel = &bargauge
		}
	case "gauge":
		var gauge GaugePanel
		p.OfType = GaugeType
		if err = json.Unmarshal(b, &gauge); err == nil {
			p.GaugePanel = &gauge
		}
	case "heatmap":
		var heatmap HeatmapPanel
		p.OfType = HeatmapType
		if err = json.Unmarshal(b, &heatmap); err == nil {
			p.HeatmapPanel = &heatmap
		}
	case "timeseries":
		var timeseries TimeseriesPanel
		p.OfType = TimeseriesType
		if err = json.Unmarshal(b, &timeseries); err == nil {
			p.TimeseriesPanel = &timeseries
		}
	case "row":
		var rowpanel RowPanel
		p.OfType = RowType
		if err = json.Unmarshal(b, &rowpanel); err == nil {
			p.RowPanel = &rowpanel
		}
	default:
		var custom = make(CustomPanel)
		p.OfType = CustomType
		if err = json.Unmarshal(b, &custom); err == nil {
			p.CustomPanel = &custom
		}
	}

	if err != nil && (probe.Title != "" || probe.Type != "") {
		err = fmt.Errorf("%w (panel %q of type %q)", err, probe.Title, probe.Type)
	}

	return err
}

func (p *Panel) MarshalJSON() ([]byte, error) {
	switch p.OfType {
	case GraphType:
		var outGraph = struct {
			CommonPanel
			GraphPanel
		}{p.CommonPanel, *p.GraphPanel}
		return json.Marshal(outGraph)
	case TableType:
		var outTable = struct {
			CommonPanel
			TablePanel
		}{p.CommonPanel, *p.TablePanel}
		return json.Marshal(outTable)
	case TextType:
		var outText = struct {
			CommonPanel
			TextPanel
		}{p.CommonPanel, *p.TextPanel}
		return json.Marshal(outText)
	case SinglestatType:
		var outSinglestat = struct {
			CommonPanel
			SinglestatPanel
		}{p.CommonPanel, *p.SinglestatPanel}
		return json.Marshal(outSinglestat)
	case StatType:
		var outSinglestat = struct {
			CommonPanel
			StatPanel
		}{p.CommonPanel, *p.StatPanel}
		return json.Marshal(outSinglestat)
	case DashlistType:
		var outDashlist = struct {
			CommonPanel
			DashlistPanel
		}{p.CommonPanel, *p.DashlistPanel}
		return json.Marshal(outDashlist)
	case BarGaugeType:
		var outBarGauge = struct {
			CommonPanel
			BarGaugePanel
		}{p.CommonPanel, *p.BarGaugePanel}
		return json.Marshal(outBarGauge)
	case GaugeType:
		var outGauge = struct {
			CommonPanel
			GaugePanel
		}{p.CommonPanel, *p.GaugePanel}
		return json.Marshal(outGauge)
	case PluginlistType:
		var outPluginlist = struct {
			CommonPanel
			PluginlistPanel
		}{p.CommonPanel, *p.PluginlistPanel}
		return json.Marshal(outPluginlist)
	case AlertlistType:
		var outAlertlist = struct {
			CommonPanel
			AlertlistPanel
		}{p.CommonPanel, *p.AlertlistPanel}
		return json.Marshal(outAlertlist)
	case RowType:
		var outRow = struct {
			CommonPanel
			RowPanel
		}{p.CommonPanel, *p.RowPanel}
		return json.Marshal(outRow)
	case HeatmapType:
		var outHeatmap = struct {
			CommonPanel
			HeatmapPanel
		}{p.CommonPanel, *p.HeatmapPanel}
		return json.Marshal(outHeatmap)
	case TimeseriesType:
		var outTimeseries = struct {
			CommonPanel
			TimeseriesPanel
		}{p.CommonPanel, *p.TimeseriesPanel}
		return json.Marshal(outTimeseries)
	case CustomType:
		var outCustom = customPanelOutput{
			p.CommonPanel,
			*p.CustomPanel,
		}
		return json.Marshal(outCustom)
	}
	return nil, errors.New("can't marshal unknown panel type")
}

type customPanelOutput struct {
	CommonPanel
	CustomPanel
}

func (c customPanelOutput) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(c.CommonPanel)
	if err != nil {
		return b, err
	}
	// Append custom keys to marshalled CommonPanel.
	buf := bytes.NewBuffer(b[:len(b)-1])

	// Sort keys to make output idempotent
	keys := make([]string, 0, len(c.CustomPanel))
	for k := range c.CustomPanel {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(`,"`)
		buf.WriteString(k)
		buf.WriteString(`":`)
		b, err := json.Marshal(c.CustomPanel[k])
		if err != nil {
			return b, err
		}
		buf.Write(b)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
