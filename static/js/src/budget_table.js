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
    const entries = this.props.entries;
    const start = formatDate(entries[0].EntryDate);
    const end = formatDate(entries[entries.length - 1].EntryDate);

    this.state = {
      startDate: start,
      endDate: end
    };
  }

  handleFilterChange = this.handleFilterChange.bind(this);

  handleFilterChange(e) {
    const value = e.target.value
    const name = e.target.name
    this.setState(state => ({
        [name]: value
    }));
  }

  render() {
    return (
      <form>
        <label className="filters">start:</label>
        <input
          type="date"
          name="startDate"
          value={this.state.startDate}
          onChange={this.handleFilterChange} />
        <br />
        <label className="filters">end:</label>
        <input
          type="date"
          name="endDate"
          value={this.state.endDate}
          onChange={this.handleFilterChange} />
        <br />
        <input type="submit" value="Submit" />
      </form>
    );
  }
}

class BudgetTable extends React.Component {
  render() {
    return (
      <div>
        <h2>budget table</h2>
        <DateFilters entries={this.props.entries}/>
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
    this.setState((state) => ({
      [name]: value
    }));
  }

  handleSubmitEntry = (e) => {
    e.preventDefault();
    // identify form values
    const entryDate = this.state.entryDate
    const amount = this.state.amount
    const category = this.state.category
    const description = this.state.description

    if ([entryDate, amount, category, description].some(i => i === '')) {
        return;
    }

    const newEntry = {
        EntryDate: this.state.entryDate,
        Amount: (this.state.amount * 100).toString(),
        Category: this.state.category,
        Description: this.state.description
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
        this.setState((state) => ({
            entryDate: '',
            amount: '',
            category: '',
            description: ''
        }));
      })
      .catch( err => console.log('something went wrong...:', err) )

    //
    console.log("constructing entry ...")


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
  constructor() {
    super();
    this.state = {
      entries: []
    };
  }

  render() {
    if (this.state.entries.length > 0) {
        return (
          <div>
            <h1>welcome!</h1>
            <EntryForm addEntry={this.handleAddEntry}/>
            <BudgetTable entries={this.state.entries}/>
          </div>
        );
    } else {
        return (<p>waiting for entries to load...</p>);
    }
  }

  handleAddEntry = (entry) => {
      this.setState( prevState => {
         return {
             entries: [
                 ...prevState.entries,
                 entry
             ]
         };
      });
  }

  componentDidMount() {
    console.log("fetching budget.json ...");
    fetch('/budget.json')
      .then(response => response.json())
      .then(responseData => {
        console.log("success!");
        console.log(responseData);
        // responseData.forEach((e) => e.Amount = e.Amount / 100);
        this.setState({ entries: responseData });
      })
      .catch(error => {
        console.log('Error fetching and parsing data', error);
      });
  }
}

const ENTRIES = [
  {EntryDate: "2021-02-01", Amount: 504.24, Category: "health", Description: "COBRA"},
  {EntryDate: "2021-02-01", Amount: 1500.00, Category: "rent", Description: "-"},
  {EntryDate: "2021-02-02", Amount: 180.85, Category: "groceries", Description: "DeCicco"},
  {EntryDate: "2021-02-03", Amount: 150.00, Category: "investing", Description: "Public.com"},
]

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
  <BudgetPage constEntries={ENTRIES}/>,
  document.querySelector('#container')
);
