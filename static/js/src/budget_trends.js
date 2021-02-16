'use strict';

// to start jsx pre-processor:
// npx babel --watch static/js/src --out-dir static/js --presets react-app/prod

class Filters extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      startDate: this.props.startDate,
      endDate: this.props.endDate,
      interval: this.props.interval,
      selectedCategories: this.props.selectedCategories,
      notYetSaved: false
    };

    this.handleSimpleChange = this.handleSimpleChange.bind(this);
    this.handleCategoryChange = this.handleCategoryChange.bind(this);
    this.createQueryString = this.createQueryString.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  render() {
    const startDate = formatDate(this.state.startDate);
    const endDate = formatDate(this.state.endDate);
    const categoryRows = []
    if (this.props.allCategories) {
      this.props.allCategories.forEach((cat, index) => {
        categoryRows.push(
          <Category
            key={index}
            name={cat} />
        );
      });
    }
    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <label>start:</label>
          <input
            type='date'
            name='startDate'
            className="filters"
            value={startDate}
            onChange={this.handleSimpleChange} />
          <br></br>
          <label>end:</label>
          <input
            type='date'
            name='endDate'
            className="filters"
            value={endDate}
            onChange={this.handleSimpleChange} />
          <br></br>
          <label>Choose categories</label>
          <select
            className='filters'
            name='selectedCategories'
            multiple={true}
            value={this.state.selectedCategories}
            size={10}
            onChange={this.handleCategoryChange} >
            {categoryRows}
          </select>
          <br></br>
          <label>Interval</label>
          <input
            type='number'
            value={this.state.interval}
            className="filters"
            name='interval'
            onChange={this.handleSimpleChange} />
          <br></br>
          <label></label>
          <input type="submit" value="Submit" />
        </form>
        {this.state.notYetSaved && <div className='edited'>edited</div>}
      </div>
    );
  }

  handleSimpleChange(event) {
    this.setState({
      [event.target.name]: event.target.value,
      notYetSaved: true
    });
  }

  handleCategoryChange(event) {
    let cats = this.state.selectedCategories;
    const value = event.target.value;
    if (cats.includes(value)) {
      cats = cats.filter(c => c != value);
    } else {
      cats.push(value);
    }
    this.setState({
      selectedCategories: cats,
      notYetSaved: true
    })
  }

  handleSubmit(event) {
    event.preventDefault();
    this.setState({ notYetSaved: false })
    const startDate = formatDate(this.state.startDate);
    const endDate = formatDate(this.state.endDate);
    const categories = this.state.selectedCategories.join('&categories=');
    const q = `?startDate=${startDate}&endDate=${endDate}&interval=${this.state.interval}&categories=${categories}`;
    this.props.fetchBudgetTrends(event, q);
  }

  createQueryString(name, value) {
    let startDate = (name == 'startDate') ? value : formatDate(this.props.startDate);
    let endDate = (name == 'endDate') ? value : formatDate(this.props.endDate);
    let interval = (name == 'interval') ? value : this.state.interval;
    let categories = this.state.selectedCategories;

    if (name == 'selectedCategories' && categories.includes(value)) {
      categories = categories.filter(c => c != value);
    } else if (name == 'selectedCategories') {
      categories.push(value);
    }
    categories = categories.join("&categories=")
    const q = `?startDate=${startDate}&endDate=${endDate}&interval=${interval}&categories=${categories}`;
    return q
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
}

class Category extends React.Component {
  render() {
    return (
      <option value={this.props.name}>
          {this.props.name}
      </option>
    );
  }
}

class HeaderRow extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    let headers = []
    if (this.props.headers) {
      headers = this.props.headers.map(h =>
        <Header key={h} name={h} />
      );
    }
    return (
      <thead>
        <tr>
          <th></th>
          {headers}
        </tr>
      </thead>
    );
  }
}

class Header extends React.Component {
  render() {
    return (
      <th>{this.props.name}</th>
    );
  }
}

class TableRows extends React.Component {
  render() {
    const rows = [];
    if (this.props.summary) {
      this.props.summary.forEach((values, i) => {
        rows.push(
          <SummaryRow
            summaryStart={this.props.dateHeaders[i]}
            values={values}
            key={i} />
        );
      });
    }

    return (
        <tbody>
          {rows}
        </tbody>
    );
  }
}

class SummaryRow extends React.Component {
  render() {
    const values = []
    this.props.values.forEach((v, i) => {
      const formattedAmount = (v / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' });
      values.push(<td key={i}>{formattedAmount}</td>);
    });

    return (
      <tr>
        <td>{this.props.summaryStart}</td>
        {values}
      </tr>
    );
  }
}

// budget over time
class BudgetTrendsContainer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      startDate: new Date(),
      endDate: new Date(),
      loaded: false
    }
  }
  render() {
    if (!this.state.loaded) {
      return null
    }
    return (
      <div>
        <p>Go to <a href="/budget-entries">budget entries</a></p>
        <h1>Budget Trends</h1>
        <Filters
          startDate={this.state.startDate}
          endDate={this.state.endDate}
          interval={this.state.interval}
          selectedCategories={this.state.selectedCategories}
          allCategories={this.state.allCategories}
          fetchBudgetTrends={this.handleFetchBudgetTrends} />
        <table>
          <HeaderRow headers={this.state.selectedCategories}/>
          <TableRows
            summary={this.state.summaryData}
            dateHeaders={this.state.dateHeaders} />
        </table>
      </div>
    );
  }

  handleFetchBudgetTrends = (e, queryString) => {
    if (e) {
      e.preventDefault();
    }
    console.log(`fetching series at: /budget-trends.json${queryString}`);
    fetch(`/budget-trends.json${queryString}`)
    .then( response => response.json() )
    .then( responseData => {
      console.log(responseData);
      this.setState (prevState => ({
        startDate: responseData['StartDate'],
        endDate: responseData['EndDate'],
        interval: responseData['TimeInterval'],
        allCategories: responseData['AllCategories'],
        selectedCategories: responseData['Table']['BucketHeaders'],
        dateHeaders: responseData['Table']['DateHeaders'],
        summaryData: responseData['Table']['Data'],
        loaded: true
      }));
    })
    .catch(error => {
      console.log('Error fetching and parsing data', error);
    });
  }

  componentDidMount() {
    this.handleFetchBudgetTrends(null, '')
  }
}

function formatDate(inputDate) {
  let options = { timeZone: 'UTC' };
  let [month, day, year] = new Date(inputDate).toLocaleDateString("en-US", options).split("/");
  if (month.length < 2) {
      month = '0' + month;
  }
  if (day.length < 2) {
      day = '0' + day;
  }
  return [year, month, day].join("-");
}

ReactDOM.render(
  <BudgetTrendsContainer />,
  document.querySelector('#container')
);
