'use strict';

// to start jsx pre-processor:
// npx babel --watch static/js/src --out-dir static/js --presets react-app/prod

class EntryForm extends React.Component {
  constructor(props) {
    super(props);
    this.handleInputChange = this.handleInputChange.bind(this)
    this.handleSubmitEntry = this.handleSubmitEntry.bind(this);
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
    console.log('submitting...')
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
      <table>
        <HeaderRow headers={this.props.headers}/>
        <TableRows
          entries={this.props.entries}
        />
      </table>
    );
  }
}

class HeaderRow extends React.Component {
  render() {
    const headers = this.props.headers.map(h =>
      <Header key={h} name={h} />
    );

    return (
      <thead>
        <tr>
          {(this.props.addCol) && <th></th>}
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
    if (this.props.entries) {
      this.props.entries.forEach((entry, index) => {
        rows.push(
          <EntryRow
            entry={entry}
            key={index} />
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

class BudgetEntriesContainer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      entries: [],
      startDate: new Date(),
      endDate: new Date()
    };
  }

  render() {
    const headers = ['EntryDate', 'Amount', 'Category', 'Description'];
    return (
      <div>
        <p>Go to <a href="/budget-trends">budget trends</a></p>
        <h1>Budget Entries</h1>
        <EntryForm addEntry={this.handleAddEntry}/>
        <DateFilters
          startDate={this.state.startDate}
          endDate={this.state.endDate}
          fetchData={this.handleFetchEntries} />
        <BudgetTable
          startDate={this.state.startDate}
          endDate={this.state.endDate}
          entries={this.state.entries}
          headers={headers}
          fetchEntries={this.handleFetchEntries}/>
      </div>
    );
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
    console.log(`fetching entries at: /budget-entries.json${queryString}`);
    fetch(`/budget-entries.json${queryString}`)
      .then(response => response.json())
      .then(responseData => {
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

// helper functions
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
  <BudgetEntriesContainer />,
  document.querySelector('#container')
);
