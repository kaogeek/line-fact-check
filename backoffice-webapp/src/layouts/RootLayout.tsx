import { Link, Outlet } from 'react-router';

export default function RootLayout() {
  return (
    <>
      <header style={{ padding: '1rem', backgroundColor: '#ddd' }}>
        <nav>
          <Link to="/">Home</Link> | <Link to="/dashboard">Dashboard</Link>
        </nav>
      </header>

      <main style={{ padding: '1rem' }}>
        <Outlet />
      </main>

      <footer style={{ padding: '1rem', backgroundColor: '#eee', marginTop: 'auto' }}>&copy; 2025 My App</footer>
    </>
  );
}
