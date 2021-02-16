'use strict';

// to start jsx pre-processor:
// npx babel --watch static/js/src --out-dir static/js --presets react-app/prod

class EntryRow extends React.Component {
  render() {
    const entry = this.props.entry
    const formattedAmount = (entry.Amount / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' });
    const formattedDate = formatDate(entry.EntryDate);
    return (
      <tr>
        <td>{formattedDate}</td>
        <td>{formattedAmount}</td>
        <td>{entry.Category}</td>
        <td>{entry.Description}</td>
      </tr>
    );
  }
}

class TableRows extends React.Component {
  render() {
    const rows = [];
    this.props.entries.forEach((entry, index) => {
      rows.push(
        <EntryRow
          entry={entry}
          key={index} />
      );
    });

    return (
        <tbody>
          {rows}
        </tbody>
    );
  }
}

class Header extends React.Component {
  render() {
    return (
      <thead>
        <tr>
          <th>Date</th>
          <th>Amount</th>
          <th>Category</th>
          <th>Description</th>
        </tr>
      </thead>
    );
  }
}

class DateFilters extends React.Component {
  constructor(props) {
    super(props);
    this.handleStartChange = this.handleStartChange.bind(this);
    this.handleEndChange = this.handleEndChange.bind(this);
  }

  render() {
    const startDate = formatDate(this.props.startDate);
    const endDate = formatDate(this.props.endDate);
    return (
      <form>
        <label className="filters">start:</label>
        <input
          type="date"
          name="startDate"
          value={startDate}
          onChange={this.handleStartChange} />
        <br />
        <label className="filters">end:</label>
        <input
          type="date"
          name="endDate"
          value={endDate}
          onChange={this.handleEndChange} />
      </form>
    );
  }

  handleStartChange(e) {
    const value = e.target.value;
    const endDate = formatDate(this.props.endDate);
    this.props.fetchData(e, `?startDate=${value}&endDate=${endDate}`)
  }

  handleEndChange(e) {
    const value = e.target.value;
    const startDate = formatDate(this.props.startDate);
    this.props.fetchData(e, `?startDate=${startDate}&endDate=${value}`)
  }
}

class BudgetTable extends React.Component {
  render() {
    return (
      <div>
        <h2>budget table</h2>
        <DateFilters
          startDate={this.props.startDate}
          endDate={this.props.endDate}
          fetchData={this.props.fetchEntries} />
        <br />
        <table>
          <Header />
          <TableRows
            entries={this.props.entries}
          />
        </table>
      </div>
    );
  }
}

class EntryForm extends React.Component {
  constructor(props) {
    super(props);
    this.handleInputChange = this.handleInputChange.bind(this)
    this.state = {
      entryDate: '',
      amount: '',
      category: '',
      description: ''
    };
  }

  render() {
    return (
      <div>
        <h2>Insert a budget entry</h2>
        <form onSubmit={this.handleSubmitEntry}>
          <label className="entry-form">entry date</label>
          <input
            name="entryDate"
            type="text"
            value={this.state.entryDate}
            onChange={this.handleInputChange} />
          <br />
          <label className="entry-form">amount</label>
          <input
            name="amount"
            type="text"
            value={this.state.amount}
            onChange={this.handleInputChange}/>
          <br />
          <label className="entry-form">category</label>
          <input
            name="category"
            type="text"
            value={this.state.category}
            onChange={this.handleInputChange} />
          <br />
          <label className="entry-form">description</label>
          <input
            name="description"
            type="text"
            value={this.state.description}
            onChange={this.handleInputChange} />
          <br />
          <input type="submit" value="Submit" />
        </form>
      </div>
    );
  }

  handleInputChange(e) {
    const name = e.target.name
    const value = e.target.value
    this.setState({ [name]: value });
  }

  handleSubmitEntry = (e) => {
    e.preventDefault();
    const entryDate = this.state.entryDate
    const amount = this.state.amount
    const category = this.state.category
    const description = this.state.description
    if ([entryDate, amount, category, description].some(i => i === '')) {
        return;
    }
    const newEntry = {
        EntryDate: entryDate,
        Amount: (amount * 100).toString(),
        Category: category,
        Description: description
    };
    // config for POST
    const config = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newEntry)
    }
    // post to db
    fetch('/insert.json', config)
      .then( response => response.json() )
      .then( responseData => {
        console.log(responseData)
        this.props.addEntry(newEntry);
        this.setState({
            entryDate: '',
            amount: '',
            category: '',
            description: ''
        });
      })
      .catch( err => console.log('something went wrong...:', err) )
  }
}

class BudgetPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      entries: []
    };
  }

  render() {
    if (this.state.entries.length == 0) {
      return null;
    } else {
      return (
        <div>
          <h1>Budget Entries</h1>
          <EntryForm addEntry={this.handleAddEntry}/>
          <BudgetTable
            startDate={this.state.startDate}
            endDate={this.state.endDate}
            entries={this.state.entries}
            fetchEntries={this.handleFetchEntries}/>
        </div>
      );
    }
  }

  handleAddEntry = (entry) => {
      this.setState( prevState => {
         return {
             entries: [...prevState.entries, entry]
         };
      });
  }

  handleFetchEntries = (e, queryString) => {
    if (e) {
      e.preventDefault();
    }
    console.log(`fetching entries at: /budget.json${queryString}`);
    fetch(`/budget.json${queryString}`)
      .then(response => response.json())
      .then(responseData => {
        console.log("... success!");
        console.log(responseData);
        this.setState( prevState => ({
          startDate: responseData['StartDate'],
          endDate: responseData['EndDate'],
          entries: responseData['Entries']
        }));
      })
      .catch(error => {
        console.log('Error fetching and parsing data', error);
      });
  }

  componentDidMount() {
    this.handleFetchEntries(null, '');
  }
}

class Category extends React.Component {
  render() {
    return (
      <option className={this.props.selectedClass} value={this.props.name}>
          {this.props.name}
      </option>
    );
  }
}

class CategoryFilter extends React.Component {
  constructor(props) {
    super(props);
    this.handleCategoryChange = this.handleCategoryChange.bind(this);
  }

  render() {
    const cats = this.props.selectedCategories;
    const categoryRows = []
    this.props.allCategories.forEach((cat, index) => {
      let selectedClass = cats.includes(cat) ? "selected" : null;
      categoryRows.push(
        <Category
          key={index}
          name={cat}
          selectedClass={selectedClass} />
      );
    });
    return (
      <form>
        <label>Choose categories</label>
        <select
          className='categories'
          multiple={true}
          value={this.props.selectedCategories}
          size={10}
          onChange={this.handleCategoryChange}>
        {categoryRows}
        </select>
      </form>
    )
  }

  handleCategoryChange(e) {
    const value = e.target.value;
    let cats = this.props.selectedCategories;
    const selected = cats.includes(value);
    if (selected) {
      cats = cats.filter(c => c != value);
    } else {
      cats.push(value);
    }
    const queryString = cats.join("&categories=");
    console.log(queryString);
    this.props.fetchData(e, `?categories=${queryString}`)
  }
}

class IntervalFilter extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
  }

  render() {
    return (
      <form>
        <label>Interval</label>
        <input type="number" value={this.props.interval}></input>
      </form>
    );
  }

  handleChange(e) {
    const value = e.target.value;
    console.log(e)
  }

}

class Filters extends React.Component {
  constructor(props) {
    super(props);
    this.handleStartChange = this.handleStartChange.bind(this);
    this.handleEndChange = this.handleEndChange.bind(this);
    this.handleCategoryChange = this.handleCategoryChange.bind(this);
    this.createQueryString = this.createQueryString.bind(this);
    // this.handleIntervalChange = this.handleIntervalChange.bind(this);

    this.handleChange = this.handleChange.bind(this);
  }

  render() {
    const startDate = formatDate(this.props.startDate);
    const endDate = formatDate(this.props.endDate);
    const cats = this.props.selectedCategories;
    const categoryRows = []
    this.props.allCategories.forEach((cat, index) => {
      let selectedClass = cats.includes(cat) ? "selected" : null;
      categoryRows.push(
        <Category
          key={index}
          name={cat}
          selectedClass={selectedClass} />
      );
    });

    return (
      <form>
        <label className="filters">start:</label>
        <input
          type='date'
          name='startDate'
          value={startDate}
          onChange={this.handleChange} />
        <label className='filters'>end:</label>
        <input
          type='date'
          name='endDate'
          value={endDate}
          onChange={this.handleChange} />
        <label>Choose categories</label>
        <select
          className='categories'
          name='categories'
          multiple={true}
          value={this.props.selectedCategories}
          size={10}
          onChange={this.handleChange} >
          {categoryRows}
        </select>
        <label>Interval</label>
        <input
          type='number'
          value={this.props.interval}
          name='interval'
          onChange={this.handleChange} />
      </form>
    );
  }

  createQueryString(name, value) {
    let startDate = (name == 'startDate') ? value : formatDate(this.props.startDate);
    let endDate = (name == 'endDate') ? value : formatDate(this.props.endDate);
    let interval = (name == 'interval') ? value : this.props.interval;
    // let selectedCategories = this.props.selectedCategories;
    let categories = this.props.selectedCategories;

    if (name == 'categories' && categories.includes(value)) {
      categories = categories.filter(c => c != value);
    } else if (name == 'categories') {
      categories.push(value);
    }
    categories = categories.join("&categories=")
    const q = `?startDate=${startDate}&endDate=${endDate}&interval=${interval}&categories=${categories}`;
    return q
  }

  handleChange(e) {
    const name = e.target.name;
    const value = e.target.value;
    console.log('name: ' + name)
    const queryString = this.createQueryString(name, value)
    console.log(queryString)
    this.props.fetchData(e, queryString)
  }

  handleStartChange(e) {
    const name = e.target.name;
    const value = e.target.value;
    const endDate = formatDate(this.props.endDate);
    this.props.fetchData(e, `?startDate=${value}&endDate=${endDate}`)
  }

  handleEndChange(e) {
    const name = e.target.name;
    const value = e.target.value;
    const startDate = formatDate(this.props.startDate);
    this.props.fetchData(e, `?startDate=${startDate}&endDate=${value}`)
  }

  handleCategoryChange(e) {
    const value = e.target.value;
    let cats = this.props.selectedCategories;
    const selected = cats.includes(value);
    if (selected) {
      cats = cats.filter(c => c != value);
    } else {
      cats.push(value);
    }
    const queryString = cats.join("&categories=");
    console.log(queryString);
    this.props.fetchData(e, `?categories=${queryString}`)
  }

}

// budget over time
class BudgetOverTimePage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: false
    }
  }
  render() {
    if (!this.state.loaded) {
      return null
    }
    return (
      <div>
      <h1>Budget Over Time</h1>
        <Filters
          startDate={this.state.startDate}
          endDate={this.state.endDate}
          interval={this.state.interval}
          selectedCategories={this.state.selectedCategories}
          allCategories={this.state.allCategories}
          fetchData={this.handleFetchBudgetOverTime} />
      </div>
    );
  }

  handleFetchBudgetOverTime = (e, queryString) => {
    if (e) {
      e.preventDefault();
    }
    console.log(`fetching series at: /budgetseries.json${queryString}`);
    fetch(`/budgetseries.json${queryString}`)
    .then( response => response.json() )
    .then( responseData => {
      console.log(responseData);
      this.setState (prevState => ({
        startDate: responseData['StartDate'],
        endDate: responseData['EndDate'],
        interval: responseData['TimeInterval'],
        allCategories: responseData['AllCategories'],
        selectedCategories: responseData['Table']['BucketHeaders'],
        queryString: queryString,
        loaded: true
      }));
    })
    .catch(error => {
      console.log('Error fetching and parsing data', error);
    });
  }

  componentDidMount() {
    this.handleFetchBudgetOverTime(null, '')
  }
}

class FinanceApp extends React.Component {
  render() {
    return (
      <div>
        <BudgetOverTimePage />
        <BudgetPage />
      </div>
    );
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
  <FinanceApp />,
  document.querySelector('#container')
);
