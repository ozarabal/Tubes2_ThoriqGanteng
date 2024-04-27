import './App.css';
import SearchBar from './components/searchBar';
import ColorToggleButton from './components/toggle';
import ColorToggleButton2 from './components/toggle2';
import React, { useState} from 'react';

function App() {
  return (
<div className="App">
  <div className="flex">
    <div className="container-fluid max-w-screen-xl px-4 mx-auto">
      <div className="col md-12">
        <div className="row my-12 border-b-4 border-r-4 border-blue-700 rounded-lg bg-blue-500 p-4">
          <h1 className="text-4xl font-bold text-center mt-8 mb-8 text-white bg-blue-500">
            Wiki Race
          </h1>
        </div>
        <div className="row my-12">
          <div className="flex flex-col items-center justify-center">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4"> {/* Grid system Tailwind */}
              <div className="w-full pr-2"> {/* Menggunakan kelas w-full untuk lebar penuh */}
                <ColorToggleButton />
              </div>
              <div className="w-full pl-2"> {/* Menggunakan kelas w-full untuk lebar penuh */}
                <ColorToggleButton2 />
              </div>
            </div>
          </div>
        </div>
        <div className="row my-12">
          <SearchBar/>
        </div>
      </div>
      </div>
  </div>
</div>

  );
}

export default App;
