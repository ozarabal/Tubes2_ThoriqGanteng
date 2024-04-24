import './App.css';
import SearchBar from './components/searchBar';
import ColorToggleButton from './components/toggle';
import ListDefault from './components/ResultBox';

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
        <div className="row my-12">
          <ColorToggleButton />
        </div>
        <div className="row my-12">
          <ListDefault />
        </div>
      </div>
    </div>
  );
}

export default App;
