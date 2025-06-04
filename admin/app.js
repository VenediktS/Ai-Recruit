const e = React.createElement;

function App() {
  const [rows, setRows] = React.useState([]);

  React.useEffect(() => {
    fetch('/admin/api/trackings')
      .then(r => r.json())
      .then(setRows)
      .catch(console.error);
  }, []);

  return e('div', null,
    e('h1', null, 'Trackings'),
    e('table', {border: 1},
      e('thead', null,
        e('tr', null,
          e('th', null, 'Компания'),
          e('th', null, 'Отправлено'),
          e('th', null, 'Переход')
        )
      ),
      e('tbody', null,
        rows.map((row, idx) => e('tr', {key: idx},
          e('td', null, row.Campaign),
          e('td', null, new Date(row.CreatedAt).toLocaleString()),
          e('td', null, row.Clicked ? 'Да' : 'Нет')
        ))
      )
    )
  );
}

ReactDOM.render(e(App), document.getElementById('root'));
