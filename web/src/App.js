import './App.css';
import SearchBar from './searchBar';

function App() {
  return (
    <div className="flex">
      <div className= "container px-8 mx-auto">
        <div className="row my-12">
          <h1 className="text-4xl font-bold text-center mt-8 mb-8">Hello World</h1>
        </div>
        <div className="row my-12">
          <SearchBar />
        </div>
      </div>
    </div>
  );
}

export default App;
