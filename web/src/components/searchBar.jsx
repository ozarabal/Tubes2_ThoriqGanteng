import React, { useState, useEffect } from 'react';

const SearchBar = ({onResponse}) => {
  const [formData, setFormData] = useState({goal: '', start: ''});
  const [response, setResponse] = useState('');

  const handleChange = (evt) => {
    const changedField = evt.target.name;
    const newValue = evt.target.value;
    setFormData( currData => {
      return {
        ...currData,
        [changedField]: newValue  
      }
    }) 
  }
  
  const handleSubmit = async (evt) => {
    evt.preventDefault();
  
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData),
    };
  
    try {
      const response = await fetch('http://localhost:8080/submit', requestOptions);
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setResponse(data);
      onResponse(data);
    } catch (error) {
      console.error('There was an error!', error);
    }
  };

  useEffect(() => {
    // Lakukan sesuatu ketika state response berubah
    console.log('State response berubah:', response);
  }, [response]); // Tambahkan response sebagai dependensi

  return (
      <div className="flex flex-col items-center justify-center">
    <form onSubmit={handleSubmit}>
      <div className="row">
        <div className="columns-md">
          <div className="row mb-4 "> 
            <label className="block text-gray-700 text-sm font-bold mb-2 text-center" htmlFor="start">
              From :
            </label>
            <input
              type="text"
              name="start"
              id="start"
              value={formData.start}
              onChange={handleChange}
              placeholder="From"
              className="border-2 border-gray-400 p-2.5 w-full rounded-md"
            />
          </div>
          <div className="row mb-4"> 
            <label className="block text-gray-700 text-sm font-bold mb-2 text-center" htmlFor="goal">
              To :
            </label>
            <input
              type="text"
              name="goal"
              id="goal"
              value={formData.goal}
              onChange={handleChange}
              placeholder="Goal"
              className="border-2 border-gray-400 p-2.5 w-full rounded-md" 
            />
          </div>
        </div>
      </div>
      <div className="row">
        <div className="flex flex-col items-center justify-center">
          <button
            type="submit"
            className="ml-4 p-4 bg-blue-500 text-white rounded-lg shadow-lg hover:bg-blue-600 focus:outline-none focus:ring mt-4"
          >
            Submit
          </button>
        </div>
      </div>
    </form>
  </div>
  );
};

export default SearchBar;
