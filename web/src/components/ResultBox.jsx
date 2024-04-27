import { List, ListItem, Card } from '@material-tailwind/react';

const ListDefault = ({ data }) => {

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

  return (
    <div className="flex flex-col items-center justify-center">
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


