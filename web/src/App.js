import './App.css';
import SearchBar from './components/searchBar';
import ColorToggleButton from './components/toggle';
import ListDefault from './components/ResultBox';
import React, { useState} from 'react';

function App() {
  const [responseData, setResponseData] = useState([]);
  
  const handleResponse  = async (data)=> {
    setResponseData(data);
  }
  return (
    <div className="flex">
      <div className= "container px-8 mx-auto">
        <div className="row my-12">
          <h1 className="text-4xl font-bold text-center mt-8 mb-8">Hello World</h1>
        </div>
        <div className="row my-12">
          <SearchBar onResponse={handleResponse}/>
        </div>
        <div className="row my-12">
          <ColorToggleButton />
        </div>
        <div className="row my-12">
          <ListDefault data={responseData}/>
        </div>
      </div>
    </div>
  );
}

export default App;
