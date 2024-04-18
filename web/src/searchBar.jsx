import React, { useState } from 'react';

const SearchBar = () => {
  const [inputValue1, setInputValue1] = useState('');
  const [inputValue2, setInputValue2] = useState('');
  const [submittedValue, setSubmittedValue] = useState(null); // State untuk menyimpan nilai yang telah disubmit

  const handleText1Change = (event) => {
    setInputValue1(event.target.value);
  };

  const handleText2Change = (event) => {
    setInputValue2(event.target.value);
  };

  const handleSubmit = () => {
    // Lakukan operasi yang diperlukan untuk mengirim data ke backend di sini
    // Misalnya, Anda dapat mengirim data melalui HTTP request menggunakan fetch atau axios
    // Di sini, kita hanya akan menyetel nilai submittedValue dengan nilai dari kedua input untuk tujuan demonstrasi
    setSubmittedValue({ input1: inputValue1, input2: inputValue2 });
  };

  return (
    <div className="flex justify-center items-center h-screen">
      <input
        type="text"
        value={inputValue1}
        onChange={handleText1Change}
        className="p-4 text-2xl border rounded-lg shadow-lg focus:outline-none focus:ring focus:border-blue-300"
        placeholder="Start"
      />
      <br />
      <input
        type="text"
        value={inputValue2}
        onChange={handleText2Change}
        className="p-4 text-2xl border rounded-lg shadow-lg focus:outline-none focus:ring focus:border-blue-300 mt-4"
        placeholder="Goal"
      />
      <br />
      <button
        onClick={handleSubmit}
        className="ml-4 p-4 bg-blue-500 text-white rounded-lg shadow-lg hover:bg-blue-600 focus:outline-none focus:ring mt-4"
      >
        Submit
      </button>
      {submittedValue && (
        <div className="mt-4">
          <p>Nilai Textbox 1 yang disubmit: {submittedValue.input1}</p>
          <p>Nilai Textbox 2 yang disubmit: {submittedValue.input2}</p>
        </div>
      )}
    </div>
  );
};

export default SearchBar;
