import React, { useState, useEffect } from 'react';
import gifImage from './thoriq.gif';
import ListDefault from './ResultBox';
import './searchBar.css';

const SearchBar = () => {
  const [formData, setFormData] = useState({ goal: '', start: '' });
  const [response, setResponse] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [search, setSearch] = useState("");
  const [suggestion, setSuggestion] = useState([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  
  const handleChange = (evt) => {
    setSearch(evt.target.value);
    setShowSuggestions(true);
    const changedField = evt.target.name;
    const newValue = evt.target.value;
    setFormData(currData => {
      return {
        ...currData,
        [changedField]: newValue
      }
    })
  }

  const handleSuggestionClick = (value, fieldName) => {
    setSearch(value);
    setFormData(prevState => ({
      ...prevState,
      [fieldName]: value
    }));
    setShowSuggestions(false);
  };
  
  const handleSubmit = async (evt) => {
    evt.preventDefault();
    setIsLoading(true);
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
      setIsLoading(false);
    } catch (error) {
      console.error('There was an error!', error);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      const limit = 5;
      if (search.trim() !== '') {
        try {
          const response = await fetch(`http://localhost:8080/fetch-wikipedia?search=${encodeURIComponent(search)}&limit=${limit}`);
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          const data = await response.json();
          console.log('Suggestion data:', data);
          setSuggestion(data);
        } catch (error) {
          console.error('Error fetching suggestions:', error);
        }
      } else {
        setSuggestion([]); // Kosongkan suggestion jika input kosong
      }
    };

    fetchData();
  }, [search]); // Jalankan efek ini setiap kali nilai search berubah

  return (
    <>
      <div className="flex flex-col items-center justify-center">
        <form onSubmit={handleSubmit}>
          <div className="row">
            <div className="columns-md flex">
              <div className="col-md-6 pr-16">
                <div className="row mb-4">
                  <label className="block text-gray-700 text-sm font-bold font-serif mb-2 text-center" htmlFor="start">
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
                  {/* Suggestions for start */}
                  {search.trim() !== '' && formData.start.trim() !== '' && showSuggestions && (
                    <div className="suggestion-container">
                      {suggestion && suggestion[1] && suggestion[1].map((item, index) => (
                        <div key={index} className="suggestion-item" onClick={() => handleSuggestionClick(item, 'start')}>{item}</div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
              <div className="col-md-6 pl-16">
                <div className="row mb-4">
                  <label className="block text-gray-700 text-sm font-bold  font-serif mb-2 text-center" htmlFor="goal">
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
                  {/* Suggestions for goal */}
                  {search.trim() !== '' && formData.goal.trim() !== '' && showSuggestions &&(
                    <div className="suggestion-container">
                      {suggestion && suggestion[1] && suggestion[1].map((item, index) => (
                        <div key={index} className="suggestion-item" onClick={() => handleSuggestionClick(item, 'goal')}>{item}</div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
          <div className="row">
            <div className="flex flex-col items-center justify-center">
              <button
                type="submit"
                className="mx-auto p-4 bg-blue-500 text-white rounded-lg shadow-lg hover:bg-blue-600 font-serif font-bold focus:outline-none focus:ring mt-4"
              >
                Submit
              </button>
            </div>
          </div>
        </form>
      </div>
      {isLoading ? <div className='flex flex-col items-center justify-center'> <img src={gifImage} alt="GIF" style={{ width: '200px', height: 'auto' }} />
</div>:<ListDefault data={response}/> }
    </>
  );
};

export default SearchBar;
