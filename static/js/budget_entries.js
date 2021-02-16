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

var EntryForm = function (_React$Component) {
  _inherits(EntryForm, _React$Component);

  function EntryForm(props) {
    _classCallCheck(this, EntryForm);

    var _this = _possibleConstructorReturn(this, (EntryForm.__proto__ || Object.getPrototypeOf(EntryForm)).call(this, props));

    _this.handleSubmitEntry = function (e) {
      e.preventDefault();
      console.log('submitting...');
      var entryDate = _this.state.entryDate;
      var amount = _this.state.amount;
      var category = _this.state.category;
      var description = _this.state.description;
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
        _this.props.addEntry(newEntry);
        _this.setState({
          entryDate: '',
          amount: '',
          category: '',
          description: ''
        });
      }).catch(function (err) {
        return console.log('something went wrong...:', err);
      });
    };

    _this.handleInputChange = _this.handleInputChange.bind(_this);
    _this.handleSubmitEntry = _this.handleSubmitEntry.bind(_this);
    _this.state = {
      entryDate: '',
      amount: '',
      category: '',
      description: ''
    };
    return _this;
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

var DateFilters = function (_React$Component2) {
  _inherits(DateFilters, _React$Component2);

  function DateFilters(props) {
    _classCallCheck(this, DateFilters);

    var _this2 = _possibleConstructorReturn(this, (DateFilters.__proto__ || Object.getPrototypeOf(DateFilters)).call(this, props));

    _this2.handleStartChange = _this2.handleStartChange.bind(_this2);
    _this2.handleEndChange = _this2.handleEndChange.bind(_this2);
    return _this2;
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

var BudgetTable = function (_React$Component3) {
  _inherits(BudgetTable, _React$Component3);

  function BudgetTable() {
    _classCallCheck(this, BudgetTable);

    return _possibleConstructorReturn(this, (BudgetTable.__proto__ || Object.getPrototypeOf(BudgetTable)).apply(this, arguments));
  }

  _createClass(BudgetTable, [{
    key: 'render',
    value: function render() {
      return React.createElement(
        'table',
        null,
        React.createElement(HeaderRow, { headers: this.props.headers }),
        React.createElement(TableRows, {
          entries: this.props.entries
        })
      );
    }
  }]);

  return BudgetTable;
}(React.Component);

var HeaderRow = function (_React$Component4) {
  _inherits(HeaderRow, _React$Component4);

  function HeaderRow() {
    _classCallCheck(this, HeaderRow);

    return _possibleConstructorReturn(this, (HeaderRow.__proto__ || Object.getPrototypeOf(HeaderRow)).apply(this, arguments));
  }

  _createClass(HeaderRow, [{
    key: 'render',
    value: function render() {
      var headers = this.props.headers.map(function (h) {
        return React.createElement(Header, { key: h, name: h });
      });

      return React.createElement(
        'thead',
        null,
        React.createElement(
          'tr',
          null,
          this.props.addCol && React.createElement('th', null),
          headers
        )
      );
    }
  }]);

  return HeaderRow;
}(React.Component);

var Header = function (_React$Component5) {
  _inherits(Header, _React$Component5);

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

var TableRows = function (_React$Component6) {
  _inherits(TableRows, _React$Component6);

  function TableRows() {
    _classCallCheck(this, TableRows);

    return _possibleConstructorReturn(this, (TableRows.__proto__ || Object.getPrototypeOf(TableRows)).apply(this, arguments));
  }

  _createClass(TableRows, [{
    key: 'render',
    value: function render() {
      var rows = [];
      if (this.props.entries) {
        this.props.entries.forEach(function (entry, index) {
          rows.push(React.createElement(EntryRow, {
            entry: entry,
            key: index }));
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

var EntryRow = function (_React$Component7) {
  _inherits(EntryRow, _React$Component7);

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

var BudgetEntriesContainer = function (_React$Component8) {
  _inherits(BudgetEntriesContainer, _React$Component8);

  function BudgetEntriesContainer(props) {
    _classCallCheck(this, BudgetEntriesContainer);

    var _this8 = _possibleConstructorReturn(this, (BudgetEntriesContainer.__proto__ || Object.getPrototypeOf(BudgetEntriesContainer)).call(this, props));

    _this8.handleAddEntry = function (entry) {
      _this8.setState(function (prevState) {
        return {
          entries: [].concat(_toConsumableArray(prevState.entries), [entry])
        };
      });
    };

    _this8.handleFetchEntries = function (e, queryString) {
      if (e) {
        e.preventDefault();
      }
      console.log('fetching entries at: /budget-entries.json' + queryString);
      fetch('/budget-entries.json' + queryString).then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log(responseData);
        _this8.setState(function (prevState) {
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

    _this8.state = {
      entries: [],
      startDate: new Date(),
      endDate: new Date()
    };
    return _this8;
  }

  _createClass(BudgetEntriesContainer, [{
    key: 'render',
    value: function render() {
      var headers = ['EntryDate', 'Amount', 'Category', 'Description'];
      return React.createElement(
        'div',
        null,
        React.createElement(
          'p',
          null,
          'Go to ',
          React.createElement(
            'a',
            { href: '/budget-trends' },
            'budget trends'
          )
        ),
        React.createElement(
          'h1',
          null,
          'Budget Entries'
        ),
        React.createElement(EntryForm, { addEntry: this.handleAddEntry }),
        React.createElement(DateFilters, {
          startDate: this.state.startDate,
          endDate: this.state.endDate,
          fetchData: this.handleFetchEntries }),
        React.createElement(BudgetTable, {
          startDate: this.state.startDate,
          endDate: this.state.endDate,
          entries: this.state.entries,
          headers: headers,
          fetchEntries: this.handleFetchEntries })
      );
    }
  }, {
    key: 'componentDidMount',
    value: function componentDidMount() {
      this.handleFetchEntries(null, '');
    }
  }]);

  return BudgetEntriesContainer;
}(React.Component);

// helper functions


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

ReactDOM.render(React.createElement(BudgetEntriesContainer, null), document.querySelector('#container'));