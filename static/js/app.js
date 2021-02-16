'use strict';

// to start jsx pre-processor:
// npx babel --watch static/js/src --out-dir static/js --presets react-app/prod

var _slicedToArray = function () { function sliceIterator(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"]) _i["return"](); } finally { if (_d) throw _e; } } return _arr; } return function (arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { return sliceIterator(arr, i); } else { throw new TypeError("Invalid attempt to destructure non-iterable instance"); } }; }();

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var EntryRow = function (_React$Component) {
  _inherits(EntryRow, _React$Component);

  function EntryRow() {
    _classCallCheck(this, EntryRow);

    return _possibleConstructorReturn(this, (EntryRow.__proto__ || Object.getPrototypeOf(EntryRow)).apply(this, arguments));
  }

  _createClass(EntryRow, [{
    key: 'render',
    value: function render() {
      var entry = this.props.entry;
      var formattedAmount = (entry.Amount / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' });
      var formattedDate = formatDate(entry.EntryDate);
      return React.createElement(
        'tr',
        null,
        React.createElement(
          'td',
          null,
          formattedDate
        ),
        React.createElement(
          'td',
          null,
          formattedAmount
        ),
        React.createElement(
          'td',
          null,
          entry.Category
        ),
        React.createElement(
          'td',
          null,
          entry.Description
        )
      );
    }
  }]);

  return EntryRow;
}(React.Component);

var TableRows = function (_React$Component2) {
  _inherits(TableRows, _React$Component2);

  function TableRows() {
    _classCallCheck(this, TableRows);

    return _possibleConstructorReturn(this, (TableRows.__proto__ || Object.getPrototypeOf(TableRows)).apply(this, arguments));
  }

  _createClass(TableRows, [{
    key: 'render',
    value: function render() {
      var rows = [];
      this.props.entries.forEach(function (entry, index) {
        rows.push(React.createElement(EntryRow, {
          entry: entry,
          key: index }));
      });

      return React.createElement(
        'tbody',
        null,
        rows
      );
    }
  }]);

  return TableRows;
}(React.Component);

var Header = function (_React$Component3) {
  _inherits(Header, _React$Component3);

  function Header() {
    _classCallCheck(this, Header);

    return _possibleConstructorReturn(this, (Header.__proto__ || Object.getPrototypeOf(Header)).apply(this, arguments));
  }

  _createClass(Header, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'thead',
        null,
        React.createElement(
          'tr',
          null,
          React.createElement(
            'th',
            null,
            'Date'
          ),
          React.createElement(
            'th',
            null,
            'Amount'
          ),
          React.createElement(
            'th',
            null,
            'Category'
          ),
          React.createElement(
            'th',
            null,
            'Description'
          )
        )
      );
    }
  }]);

  return Header;
}(React.Component);

var DateFilters = function (_React$Component4) {
  _inherits(DateFilters, _React$Component4);

  function DateFilters(props) {
    _classCallCheck(this, DateFilters);

    var _this4 = _possibleConstructorReturn(this, (DateFilters.__proto__ || Object.getPrototypeOf(DateFilters)).call(this, props));

    _this4.handleStartChange = _this4.handleStartChange.bind(_this4);
    _this4.handleEndChange = _this4.handleEndChange.bind(_this4);
    return _this4;
  }

  _createClass(DateFilters, [{
    key: 'render',
    value: function render() {
      var startDate = formatDate(this.props.startDate);
      var endDate = formatDate(this.props.endDate);
      return React.createElement(
        'form',
        null,
        React.createElement(
          'label',
          { className: 'filters' },
          'start:'
        ),
        React.createElement('input', {
          type: 'date',
          name: 'startDate',
          value: startDate,
          onChange: this.handleStartChange }),
        React.createElement('br', null),
        React.createElement(
          'label',
          { className: 'filters' },
          'end:'
        ),
        React.createElement('input', {
          type: 'date',
          name: 'endDate',
          value: endDate,
          onChange: this.handleEndChange })
      );
    }
  }, {
    key: 'handleStartChange',
    value: function handleStartChange(e) {
      var value = e.target.value;
      var endDate = formatDate(this.props.endDate);
      this.props.fetchData(e, '?startDate=' + value + '&endDate=' + endDate);
    }
  }, {
    key: 'handleEndChange',
    value: function handleEndChange(e) {
      var value = e.target.value;
      var startDate = formatDate(this.props.startDate);
      this.props.fetchData(e, '?startDate=' + startDate + '&endDate=' + value);
    }
  }]);

  return DateFilters;
}(React.Component);

var BudgetTable = function (_React$Component5) {
  _inherits(BudgetTable, _React$Component5);

  function BudgetTable() {
    _classCallCheck(this, BudgetTable);

    return _possibleConstructorReturn(this, (BudgetTable.__proto__ || Object.getPrototypeOf(BudgetTable)).apply(this, arguments));
  }

  _createClass(BudgetTable, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'div',
        null,
        React.createElement(
          'h2',
          null,
          'budget table'
        ),
        React.createElement(DateFilters, {
          startDate: this.props.startDate,
          endDate: this.props.endDate,
          fetchData: this.props.fetchEntries }),
        React.createElement('br', null),
        React.createElement(
          'table',
          null,
          React.createElement(Header, null),
          React.createElement(TableRows, {
            entries: this.props.entries
          })
        )
      );
    }
  }]);

  return BudgetTable;
}(React.Component);

var EntryForm = function (_React$Component6) {
  _inherits(EntryForm, _React$Component6);

  function EntryForm(props) {
    _classCallCheck(this, EntryForm);

    var _this6 = _possibleConstructorReturn(this, (EntryForm.__proto__ || Object.getPrototypeOf(EntryForm)).call(this, props));

    _this6.handleSubmitEntry = function (e) {
      e.preventDefault();
      var entryDate = _this6.state.entryDate;
      var amount = _this6.state.amount;
      var category = _this6.state.category;
      var description = _this6.state.description;
      if ([entryDate, amount, category, description].some(function (i) {
        return i === '';
      })) {
        return;
      }
      var newEntry = {
        EntryDate: entryDate,
        Amount: (amount * 100).toString(),
        Category: category,
        Description: description
      };
      // config for POST
      var config = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newEntry)
        // post to db
      };fetch('/insert.json', config).then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log(responseData);
        _this6.props.addEntry(newEntry);
        _this6.setState({
          entryDate: '',
          amount: '',
          category: '',
          description: ''
        });
      }).catch(function (err) {
        return console.log('something went wrong...:', err);
      });
    };

    _this6.handleInputChange = _this6.handleInputChange.bind(_this6);
    _this6.state = {
      entryDate: '',
      amount: '',
      category: '',
      description: ''
    };
    return _this6;
  }

  _createClass(EntryForm, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'div',
        null,
        React.createElement(
          'h2',
          null,
          'Insert a budget entry'
        ),
        React.createElement(
          'form',
          { onSubmit: this.handleSubmitEntry },
          React.createElement(
            'label',
            { className: 'entry-form' },
            'entry date'
          ),
          React.createElement('input', {
            name: 'entryDate',
            type: 'text',
            value: this.state.entryDate,
            onChange: this.handleInputChange }),
          React.createElement('br', null),
          React.createElement(
            'label',
            { className: 'entry-form' },
            'amount'
          ),
          React.createElement('input', {
            name: 'amount',
            type: 'text',
            value: this.state.amount,
            onChange: this.handleInputChange }),
          React.createElement('br', null),
          React.createElement(
            'label',
            { className: 'entry-form' },
            'category'
          ),
          React.createElement('input', {
            name: 'category',
            type: 'text',
            value: this.state.category,
            onChange: this.handleInputChange }),
          React.createElement('br', null),
          React.createElement(
            'label',
            { className: 'entry-form' },
            'description'
          ),
          React.createElement('input', {
            name: 'description',
            type: 'text',
            value: this.state.description,
            onChange: this.handleInputChange }),
          React.createElement('br', null),
          React.createElement('input', { type: 'submit', value: 'Submit' })
        )
      );
    }
  }, {
    key: 'handleInputChange',
    value: function handleInputChange(e) {
      var name = e.target.name;
      var value = e.target.value;
      this.setState(_defineProperty({}, name, value));
    }
  }]);

  return EntryForm;
}(React.Component);

var BudgetPage = function (_React$Component7) {
  _inherits(BudgetPage, _React$Component7);

  function BudgetPage(props) {
    _classCallCheck(this, BudgetPage);

    var _this7 = _possibleConstructorReturn(this, (BudgetPage.__proto__ || Object.getPrototypeOf(BudgetPage)).call(this, props));

    _this7.handleAddEntry = function (entry) {
      _this7.setState(function (prevState) {
        return {
          entries: [].concat(_toConsumableArray(prevState.entries), [entry])
        };
      });
    };

    _this7.handleFetchEntries = function (e, queryString) {
      if (e) {
        e.preventDefault();
      }
      console.log('fetching entries at: /budget.json' + queryString);
      fetch('/budget.json' + queryString).then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log("... success!");
        console.log(responseData);
        _this7.setState(function (prevState) {
          return {
            startDate: responseData['StartDate'],
            endDate: responseData['EndDate'],
            entries: responseData['Entries']
          };
        });
      }).catch(function (error) {
        console.log('Error fetching and parsing data', error);
      });
    };

    _this7.state = {
      entries: []
    };
    return _this7;
  }

  _createClass(BudgetPage, [{
    key: 'render',
    value: function render() {
      if (this.state.entries.length == 0) {
        return null;
      } else {
        return React.createElement(
          'div',
          null,
          React.createElement(
            'h1',
            null,
            'Budget Entries'
          ),
          React.createElement(EntryForm, { addEntry: this.handleAddEntry }),
          React.createElement(BudgetTable, {
            startDate: this.state.startDate,
            endDate: this.state.endDate,
            entries: this.state.entries,
            fetchEntries: this.handleFetchEntries })
        );
      }
    }
  }, {
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.handleFetchEntries(null, '');
    }
  }]);

  return BudgetPage;
}(React.Component);

var Category = function (_React$Component8) {
  _inherits(Category, _React$Component8);

  function Category() {
    _classCallCheck(this, Category);

    return _possibleConstructorReturn(this, (Category.__proto__ || Object.getPrototypeOf(Category)).apply(this, arguments));
  }

  _createClass(Category, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'option',
        { className: this.props.selectedClass, value: this.props.name },
        this.props.name
      );
    }
  }]);

  return Category;
}(React.Component);

var CategoryFilter = function (_React$Component9) {
  _inherits(CategoryFilter, _React$Component9);

  function CategoryFilter(props) {
    _classCallCheck(this, CategoryFilter);

    var _this9 = _possibleConstructorReturn(this, (CategoryFilter.__proto__ || Object.getPrototypeOf(CategoryFilter)).call(this, props));

    _this9.handleCategoryChange = _this9.handleCategoryChange.bind(_this9);
    return _this9;
  }

  _createClass(CategoryFilter, [{
    key: 'render',
    value: function render() {
      var cats = this.props.selectedCategories;
      var categoryRows = [];
      this.props.allCategories.forEach(function (cat, index) {
        var selectedClass = cats.includes(cat) ? "selected" : null;
        categoryRows.push(React.createElement(Category, {
          key: index,
          name: cat,
          selectedClass: selectedClass }));
      });
      return React.createElement(
        'form',
        null,
        React.createElement(
          'label',
          null,
          'Choose categories'
        ),
        React.createElement(
          'select',
          {
            className: 'categories',
            multiple: true,
            value: this.props.selectedCategories,
            size: 10,
            onChange: this.handleCategoryChange },
          categoryRows
        )
      );
    }
  }, {
    key: 'handleCategoryChange',
    value: function handleCategoryChange(e) {
      var value = e.target.value;
      var cats = this.props.selectedCategories;
      var selected = cats.includes(value);
      if (selected) {
        cats = cats.filter(function (c) {
          return c != value;
        });
      } else {
        cats.push(value);
      }
      var queryString = cats.join("&categories=");
      console.log(queryString);
      this.props.fetchData(e, '?categories=' + queryString);
    }
  }]);

  return CategoryFilter;
}(React.Component);

var IntervalFilter = function (_React$Component10) {
  _inherits(IntervalFilter, _React$Component10);

  function IntervalFilter(props) {
    _classCallCheck(this, IntervalFilter);

    var _this10 = _possibleConstructorReturn(this, (IntervalFilter.__proto__ || Object.getPrototypeOf(IntervalFilter)).call(this, props));

    _this10.handleChange = _this10.handleChange.bind(_this10);
    return _this10;
  }

  _createClass(IntervalFilter, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'form',
        null,
        React.createElement(
          'label',
          null,
          'Interval'
        ),
        React.createElement('input', { type: 'number', value: this.props.interval })
      );
    }
  }, {
    key: 'handleChange',
    value: function handleChange(e) {
      var value = e.target.value;
      console.log(e);
    }
  }]);

  return IntervalFilter;
}(React.Component);

var Filters = function (_React$Component11) {
  _inherits(Filters, _React$Component11);

  function Filters() {
    _classCallCheck(this, Filters);

    return _possibleConstructorReturn(this, (Filters.__proto__ || Object.getPrototypeOf(Filters)).apply(this, arguments));
  }

  _createClass(Filters, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'div',
        null,
        React.createElement(DateFilters, {
          startDate: this.props.startDate,
          endDate: this.props.endDate,
          fetchData: this.props.fetchData }),
        React.createElement(CategoryFilter, {
          selectedCategories: this.props.selectedCategories,
          allCategories: this.props.allCategories,
          fetchData: this.props.fetchData }),
        React.createElement(IntervalFilter, null)
      );
    }
  }]);

  return Filters;
}(React.Component);

// budget over time


var BudgetOverTimePage = function (_React$Component12) {
  _inherits(BudgetOverTimePage, _React$Component12);

  function BudgetOverTimePage(props) {
    _classCallCheck(this, BudgetOverTimePage);

    var _this12 = _possibleConstructorReturn(this, (BudgetOverTimePage.__proto__ || Object.getPrototypeOf(BudgetOverTimePage)).call(this, props));

    _this12.handleFetchBudgetOverTime = function (e, queryString) {
      if (e) {
        e.preventDefault();
      }
      console.log('fetching series at: /budgetseries.json' + queryString);
      fetch('/budgetseries.json' + queryString).then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log(responseData);
        _this12.setState(function (prevState) {
          return {
            startDate: responseData['StartDate'],
            endDate: responseData['EndDate'],
            allCategories: responseData['AllCategories'],
            selectedCategories: responseData['Table']['BucketHeaders'],
            queryString: queryString,
            loaded: true
          };
        });
      }).catch(function (error) {
        console.log('Error fetching and parsing data', error);
      });
    };

    _this12.state = {
      loaded: false
    };
    return _this12;
  }

  _createClass(BudgetOverTimePage, [{
    key: 'render',
    value: function render() {
      if (!this.state.loaded) {
        return null;
      }
      return React.createElement(
        'div',
        null,
        React.createElement(
          'h1',
          null,
          'Budget Over Time'
        ),
        React.createElement(Filters, {
          startDate: this.state.startDate,
          endDate: this.state.endDate,
          selectedCategories: this.state.selectedCategories,
          allCategories: this.state.allCategories,
          fetchData: this.handleFetchBudgetOverTime })
      );
    }
  }, {
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.handleFetchBudgetOverTime(null, '');
    }
  }]);

  return BudgetOverTimePage;
}(React.Component);

var FinanceApp = function (_React$Component13) {
  _inherits(FinanceApp, _React$Component13);

  function FinanceApp() {
    _classCallCheck(this, FinanceApp);

    return _possibleConstructorReturn(this, (FinanceApp.__proto__ || Object.getPrototypeOf(FinanceApp)).apply(this, arguments));
  }

  _createClass(FinanceApp, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'div',
        null,
        React.createElement(BudgetOverTimePage, null),
        React.createElement(BudgetPage, null)
      );
    }
  }]);

  return FinanceApp;
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

ReactDOM.render(React.createElement(FinanceApp, null), document.querySelector('#container'));