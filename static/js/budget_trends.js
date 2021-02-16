'use strict';

// to start jsx pre-processor:
// npx babel --watch static/js/src --out-dir static/js --presets react-app/prod

var _slicedToArray = function () { function sliceIterator(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"]) _i["return"](); } finally { if (_d) throw _e; } } return _arr; } return function (arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { return sliceIterator(arr, i); } else { throw new TypeError("Invalid attempt to destructure non-iterable instance"); } }; }();

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var Filters = function (_React$Component) {
  _inherits(Filters, _React$Component);

  function Filters(props) {
    _classCallCheck(this, Filters);

    var _this = _possibleConstructorReturn(this, (Filters.__proto__ || Object.getPrototypeOf(Filters)).call(this, props));

    _this.state = {
      startDate: _this.props.startDate,
      endDate: _this.props.endDate,
      interval: _this.props.interval,
      selectedCategories: _this.props.selectedCategories,
      notYetSaved: false
    };

    _this.handleSimpleChange = _this.handleSimpleChange.bind(_this);
    _this.handleCategoryChange = _this.handleCategoryChange.bind(_this);
    _this.createQueryString = _this.createQueryString.bind(_this);
    _this.handleSubmit = _this.handleSubmit.bind(_this);
    return _this;
  }

  _createClass(Filters, [{
    key: 'render',
    value: function render() {
      var startDate = formatDate(this.state.startDate);
      var endDate = formatDate(this.state.endDate);
      var categoryRows = [];
      if (this.props.allCategories) {
        this.props.allCategories.forEach(function (cat, index) {
          categoryRows.push(React.createElement(Category, {
            key: index,
            name: cat }));
        });
      }
      return React.createElement(
        'div',
        null,
        React.createElement(
          'form',
          { onSubmit: this.handleSubmit },
          React.createElement(
            'label',
            null,
            'start:'
          ),
          React.createElement('input', {
            type: 'date',
            name: 'startDate',
            className: 'filters',
            value: startDate,
            onChange: this.handleSimpleChange }),
          React.createElement('br', null),
          React.createElement(
            'label',
            null,
            'end:'
          ),
          React.createElement('input', {
            type: 'date',
            name: 'endDate',
            className: 'filters',
            value: endDate,
            onChange: this.handleSimpleChange }),
          React.createElement('br', null),
          React.createElement(
            'label',
            null,
            'Choose categories'
          ),
          React.createElement(
            'select',
            {
              className: 'filters',
              name: 'selectedCategories',
              multiple: true,
              value: this.state.selectedCategories,
              size: 10,
              onChange: this.handleCategoryChange },
            categoryRows
          ),
          React.createElement('br', null),
          React.createElement(
            'label',
            null,
            'Interval'
          ),
          React.createElement('input', {
            type: 'number',
            value: this.state.interval,
            className: 'filters',
            name: 'interval',
            onChange: this.handleSimpleChange }),
          React.createElement('br', null),
          React.createElement('label', null),
          React.createElement('input', { type: 'submit', value: 'Submit' })
        ),
        this.state.notYetSaved && React.createElement(
          'div',
          { className: 'edited' },
          'edited'
        )
      );
    }
  }, {
    key: 'handleSimpleChange',
    value: function handleSimpleChange(event) {
      var _setState;

      this.setState((_setState = {}, _defineProperty(_setState, event.target.name, event.target.value), _defineProperty(_setState, 'notYetSaved', true), _setState));
    }
  }, {
    key: 'handleCategoryChange',
    value: function handleCategoryChange(event) {
      var cats = this.state.selectedCategories;
      var value = event.target.value;
      if (cats.includes(value)) {
        cats = cats.filter(function (c) {
          return c != value;
        });
      } else {
        cats.push(value);
      }
      this.setState({
        selectedCategories: cats,
        notYetSaved: true
      });
    }
  }, {
    key: 'handleSubmit',
    value: function handleSubmit(event) {
      event.preventDefault();
      this.setState({ notYetSaved: false });
      var startDate = formatDate(this.state.startDate);
      var endDate = formatDate(this.state.endDate);
      var categories = this.state.selectedCategories.join('&categories=');
      var q = '?startDate=' + startDate + '&endDate=' + endDate + '&interval=' + this.state.interval + '&categories=' + categories;
      this.props.fetchBudgetTrends(event, q);
    }
  }, {
    key: 'createQueryString',
    value: function createQueryString(name, value) {
      var startDate = name == 'startDate' ? value : formatDate(this.props.startDate);
      var endDate = name == 'endDate' ? value : formatDate(this.props.endDate);
      var interval = name == 'interval' ? value : this.state.interval;
      var categories = this.state.selectedCategories;

      if (name == 'selectedCategories' && categories.includes(value)) {
        categories = categories.filter(function (c) {
          return c != value;
        });
      } else if (name == 'selectedCategories') {
        categories.push(value);
      }
      categories = categories.join("&categories=");
      var q = '?startDate=' + startDate + '&endDate=' + endDate + '&interval=' + interval + '&categories=' + categories;
      return q;
    }

    // handleStartChange(e) {
    //   const value = e.target.value;
    //   const endDate = formatDate(this.props.endDate);
    //   const q = this.createQueryString('startDate', value);
    //   console.log(q);
    //   this.props.fetchData(e, q)
    // }
    //
    // handleEndChange(e) {
    //   const value = e.target.value;
    //   const startDate = formatDate(this.props.startDate);
    //   this.props.fetchData(e, `?startDate=${startDate}&endDate=${value}`)
    // }

  }]);

  return Filters;
}(React.Component);

var Category = function (_React$Component2) {
  _inherits(Category, _React$Component2);

  function Category() {
    _classCallCheck(this, Category);

    return _possibleConstructorReturn(this, (Category.__proto__ || Object.getPrototypeOf(Category)).apply(this, arguments));
  }

  _createClass(Category, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'option',
        { value: this.props.name },
        this.props.name
      );
    }
  }]);

  return Category;
}(React.Component);

var HeaderRow = function (_React$Component3) {
  _inherits(HeaderRow, _React$Component3);

  function HeaderRow(props) {
    _classCallCheck(this, HeaderRow);

    return _possibleConstructorReturn(this, (HeaderRow.__proto__ || Object.getPrototypeOf(HeaderRow)).call(this, props));
  }

  _createClass(HeaderRow, [{
    key: 'render',
    value: function render() {
      var headers = [];
      if (this.props.headers) {
        headers = this.props.headers.map(function (h) {
          return React.createElement(Header, { key: h, name: h });
        });
      }
      return React.createElement(
        'thead',
        null,
        React.createElement(
          'tr',
          null,
          React.createElement('th', null),
          headers
        )
      );
    }
  }]);

  return HeaderRow;
}(React.Component);

var Header = function (_React$Component4) {
  _inherits(Header, _React$Component4);

  function Header() {
    _classCallCheck(this, Header);

    return _possibleConstructorReturn(this, (Header.__proto__ || Object.getPrototypeOf(Header)).apply(this, arguments));
  }

  _createClass(Header, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'th',
        null,
        this.props.name
      );
    }
  }]);

  return Header;
}(React.Component);

var TableRows = function (_React$Component5) {
  _inherits(TableRows, _React$Component5);

  function TableRows() {
    _classCallCheck(this, TableRows);

    return _possibleConstructorReturn(this, (TableRows.__proto__ || Object.getPrototypeOf(TableRows)).apply(this, arguments));
  }

  _createClass(TableRows, [{
    key: 'render',
    value: function render() {
      var _this6 = this;

      var rows = [];
      if (this.props.summary) {
        this.props.summary.forEach(function (values, i) {
          rows.push(React.createElement(SummaryRow, {
            summaryStart: _this6.props.dateHeaders[i],
            values: values,
            key: i }));
        });
      }

      return React.createElement(
        'tbody',
        null,
        rows
      );
    }
  }]);

  return TableRows;
}(React.Component);

var SummaryRow = function (_React$Component6) {
  _inherits(SummaryRow, _React$Component6);

  function SummaryRow() {
    _classCallCheck(this, SummaryRow);

    return _possibleConstructorReturn(this, (SummaryRow.__proto__ || Object.getPrototypeOf(SummaryRow)).apply(this, arguments));
  }

  _createClass(SummaryRow, [{
    key: 'render',
    value: function render() {
      var values = [];
      this.props.values.forEach(function (v, i) {
        var formattedAmount = (v / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' });
        values.push(React.createElement(
          'td',
          { key: i },
          formattedAmount
        ));
      });

      return React.createElement(
        'tr',
        null,
        React.createElement(
          'td',
          null,
          this.props.summaryStart
        ),
        values
      );
    }
  }]);

  return SummaryRow;
}(React.Component);

// budget over time


var BudgetTrendsContainer = function (_React$Component7) {
  _inherits(BudgetTrendsContainer, _React$Component7);

  function BudgetTrendsContainer(props) {
    _classCallCheck(this, BudgetTrendsContainer);

    var _this8 = _possibleConstructorReturn(this, (BudgetTrendsContainer.__proto__ || Object.getPrototypeOf(BudgetTrendsContainer)).call(this, props));

    _this8.handleFetchBudgetTrends = function (e, queryString) {
      if (e) {
        e.preventDefault();
      }
      console.log('fetching series at: /budget-trends.json' + queryString);
      fetch('/budget-trends.json' + queryString).then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log(responseData);
        _this8.setState(function (prevState) {
          return {
            startDate: responseData['StartDate'],
            endDate: responseData['EndDate'],
            interval: responseData['TimeInterval'],
            allCategories: responseData['AllCategories'],
            selectedCategories: responseData['Table']['BucketHeaders'],
            dateHeaders: responseData['Table']['DateHeaders'],
            summaryData: responseData['Table']['Data'],
            loaded: true
          };
        });
      }).catch(function (error) {
        console.log('Error fetching and parsing data', error);
      });
    };

    _this8.state = {
      startDate: new Date(),
      endDate: new Date(),
      loaded: false
    };
    return _this8;
  }

  _createClass(BudgetTrendsContainer, [{
    key: 'render',
    value: function render() {
      if (!this.state.loaded) {
        return null;
      }
      return React.createElement(
        'div',
        null,
        React.createElement(
          'p',
          null,
          'Go to ',
          React.createElement(
            'a',
            { href: '/budget-entries' },
            'budget entries'
          )
        ),
        React.createElement(
          'h1',
          null,
          'Budget Trends'
        ),
        React.createElement(Filters, {
          startDate: this.state.startDate,
          endDate: this.state.endDate,
          interval: this.state.interval,
          selectedCategories: this.state.selectedCategories,
          allCategories: this.state.allCategories,
          fetchBudgetTrends: this.handleFetchBudgetTrends }),
        React.createElement(
          'table',
          null,
          React.createElement(HeaderRow, { headers: this.state.selectedCategories }),
          React.createElement(TableRows, {
            summary: this.state.summaryData,
            dateHeaders: this.state.dateHeaders })
        )
      );
    }
  }, {
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.handleFetchBudgetTrends(null, '');
    }
  }]);

  return BudgetTrendsContainer;
}(React.Component);

function formatDate(inputDate) {
  var options = { timeZone: 'UTC' };

  var _toLocaleDateString$s = new Date(inputDate).toLocaleDateString("en-US", options).split("/"),
      _toLocaleDateString$s2 = _slicedToArray(_toLocaleDateString$s, 3),
      month = _toLocaleDateString$s2[0],
      day = _toLocaleDateString$s2[1],
      year = _toLocaleDateString$s2[2];

  if (month.length < 2) {
    month = '0' + month;
  }
  if (day.length < 2) {
    day = '0' + day;
  }
  return [year, month, day].join("-");
}

ReactDOM.render(React.createElement(BudgetTrendsContainer, null), document.querySelector('#container'));