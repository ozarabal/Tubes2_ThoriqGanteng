import { List, ListItem, Card } from '@material-tailwind/react';

const ListDefault = ({ data }) => {
  console.log(data);
  if (!data || !data.result || !Array.isArray(data.result) || data.result.length === 0) {
    return null; // Kembalikan null jika tidak ada data atau data tidak sesuai
  }

  // Ubah data menjadi format yang diinginkan dengan menambahkan spasi setiap i+1

  
  const formattedData = data.result.map((row, i) =>
  row.map((item, j) => {
      const linkText = item.replace('https://en.wikipedia.org/wiki/', ''); // Menghapus bagian yang tidak diinginkan dari tautan
      return (
        <a key={`${i}-${j}`} href={item} target="_blank" rel="noopener noreferrer">
          {linkText}
        </a>
      );
    })
  );

  const spacedFormattedData = formattedData.map((row, i) => (
    <div key={`row-${i}`}>
      {row.map((link, j) => (
        <span key={`link-${i}-${j}`}>
          {link}
          {j < row.length - 1 && ' '} {/* Tambahkan spasi jika bukan elemen terakhir dalam baris */}
        </span>
      ))}
      {i < formattedData.length - 1 && <br />} {/* Tambahkan baris baru jika bukan baris terakhir */}
    </div>
  ));
  const banyak_path = data.banyak_path;
  const waktup = data.waktu;
  const banyaklink = data.banyak_jelajah;
  const kedalaman = data.kedalaman;

  return (
    <div className="flex flex-col items-center justify-center">
      <div className='row mb-10'>
        <div className='col flex border-b-4 border-r-4 border-blue-700 rounded-lg bg-blue-500 p-4'>
          <div className='col md-6 mr-8'>
            <div className='row'>
              <h1 className="text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">Banyak Path:</h1>
              <h2 className= "text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">{banyak_path}</h2>
            </div>
            <div className='row'>
              <h1 className="text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">Banyaknya link :</h1>
                <h2 className= "text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">{banyaklink}</h2>
            </div> 
          </div>
          <div className='col md-6 ml-8'>
            <div className='row'>
              <h1 className="text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">Waktu Pencarian :</h1>
              <h2 className="text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">{waktup} detik</h2>
            </div>
            <div className='row'>
              <h1 className="text-2x1 font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">Artikel dilalui :</h1>
              <h2 className="text-2xl font-bold font-serif text-center mt-8 mb-8 text-white bg-blue-500">{kedalaman}</h2>
            </div>
          </div>
        </div>
      </div>
      <Card className="w-200">
        <List>
          <ListItem>
            <pre>{spacedFormattedData}</pre>
          </ListItem>
        </List>
      </Card>
    </div>
  );
};

export default ListDefault;


