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
          key={index}/>
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
  }

  handleStartChange = this.handleStartChange.bind(this);
  handleEndChange = this.handleEndChange.bind(this);

  handleStartChange(e) {
    const value = e.target.value;
    const endDate = formatDate(this.props.endDate);
    this.props.fetchEntries(e, `/budget.json?startDate=${value}&endDate=${endDate}`)
  }

  handleEndChange(e) {
    const value = e.target.value;
    const startDate = formatDate(this.props.startDate);
    this.props.fetchEntries(e, `/budget.json?startDate=${startDate}&endDate=${value}`)
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
        <br />
        <input type="submit" value="Submit" />
      </form>
    );
  }
}

class BudgetTable extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>
        <h2>budget table</h2>
        <DateFilters
          startDate={this.props.startDate}
          endDate={this.props.endDate}
          fetchEntries={this.props.fetchEntries}/>
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
    this.state = {
      entryDate: '',
      amount: '',
      category: '',
      description: ''
    };
  }

  handleInputChange = this.handleInputChange.bind(this)

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
          <h1>welcome!</h1>
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

  handleFetchEntries = (e, fetchUrl) => {
    e.preventDefault();
    console.log(`fetching entries ... fetch url: ${fetchUrl}`);
    fetch(fetchUrl)
      .then(response => response.json())
      .then(responseData => {
        console.log("... success!");
        console.log(responseData);
        this.setState( prevState => ({
          startDate: responseData['Start'],
          endDate: responseData['End'],
          entries: responseData['Entries']
        }));
      })
      .catch(error => {
        console.log('Error fetching and parsing data', error);
      });
  }

  componentDidMount() {
    console.log("componentDidMount. this.state.startDate:" + this.state.startDate)
    fetch('/budget.json')
      .then(response => response.json())
      .then(responseData => {
        console.log("... success!");
        console.log(responseData);
        this.setState( prevState => ({
          startDate: responseData['Start'],
          endDate: responseData['End'],
          entries: responseData['Entries']
        }));
      })
      .catch(error => {
        console.log('Error fetching and parsing data', error);
      });
  }
}

function formatDate(inputDate) {
  let [month, day, year] = new Date(inputDate).toLocaleDateString("en-US").split("/");
  if (month.length < 2) {
      month = '0' + month;
  }
  if (day.length < 2) {
      day = '0' + day;
  }
  return [year, month, day].join("-");
}

ReactDOM.render(
  <BudgetPage />,
  document.querySelector('#container')
);
