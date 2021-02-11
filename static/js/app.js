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

    _this4.handleFilterChange = _this4.handleFilterChange.bind(_this4);

    _this4.handleSubmitDateFilters = function (e) {
      console.log("inside handleSubmitDateFilters: " + _this4.state.startDate);
      _this4.props.fetchEntries(e, '/budget.json?start=' + _this4.state.startDate);
    };

    var start = formatDate(_this4.props.startDate);
    var end = formatDate(_this4.props.endDate);

    _this4.state = {
      startDate: start,
      endDate: end
    };
    return _this4;
  }

  _createClass(DateFilters, [{
    key: 'handleFilterChange',
    value: function handleFilterChange(e) {
      var value = e.target.value;
      var name = e.target.name;
      this.setState(_defineProperty({}, name, value));
    }
  }, {
    key: 'render',
    value: function render() {
      return React.createElement(
        'form',
        { onSubmit: this.handleSubmitDateFilters },
        React.createElement(
          'label',
          { className: 'filters' },
          'start:'
        ),
        React.createElement('input', {
          type: 'date',
          name: 'startDate',
          value: this.state.startDate,
          onChange: this.handleFilterChange }),
        React.createElement('br', null),
        React.createElement(
          'label',
          { className: 'filters' },
          'end:'
        ),
        React.createElement('input', {
          type: 'date',
          name: 'endDate',
          value: this.state.endDate,
          onChange: this.handleFilterChange }),
        React.createElement('br', null),
        React.createElement('input', { type: 'submit', value: 'Submit' })
      );
    }
  }]);

  return DateFilters;
}(React.Component);

var BudgetTable = function (_React$Component5) {
  _inherits(BudgetTable, _React$Component5);

  function BudgetTable(props) {
    _classCallCheck(this, BudgetTable);

    return _possibleConstructorReturn(this, (BudgetTable.__proto__ || Object.getPrototypeOf(BudgetTable)).call(this, props));
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
          fetchEntries: this.props.fetchEntries }),
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

    _this6.handleInputChange = _this6.handleInputChange.bind(_this6);

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

    _this6.state = {
      entryDate: '',
      amount: '',
      category: '',
      description: ''
    };
    return _this6;
  }

  _createClass(EntryForm, [{
    key: 'handleInputChange',
    value: function handleInputChange(e) {
      var name = e.target.name;
      var value = e.target.value;
      this.setState(_defineProperty({}, name, value));
    }
  }, {
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

    _this7.handleFetchEntries = function (e, fetchUrl) {
      e.preventDefault();
      console.log('fetching entries ... fetch url: ' + fetchUrl);
      fetch(fetchUrl).then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log("... success!");
        console.log(responseData);
        _this7.setState(function (prevState) {
          return {
            startDate: responseData['Start'],
            endDate: responseData['End'],
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
            'welcome!'
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
      var _this8 = this;

      console.log("componentDidMount. this.state.startDate:" + this.state.startDate);
      fetch('/budget.json').then(function (response) {
        return response.json();
      }).then(function (responseData) {
        console.log("... success!");
        console.log(responseData);
        _this8.setState(function (prevState) {
          return {
            startDate: responseData['Start'],
            endDate: responseData['End'],
            entries: responseData['Entries']
          };
        });
      }).catch(function (error) {
        console.log('Error fetching and parsing data', error);
      });
    }
  }]);

  return BudgetPage;
}(React.Component);

function formatDate(inputDate) {
  var _toLocaleDateString$s = new Date(inputDate).toLocaleDateString("en-US").split("/"),
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

ReactDOM.render(React.createElement(BudgetPage, null), document.querySelector('#container'));