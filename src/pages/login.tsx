import React, { useState, type ChangeEvent, type FormEvent } from 'react';
import { Alert, Button, Container, Form } from 'react-bootstrap';
import { useNavigate } from "react-router";
import { ROUTES } from '../routes';
import { loginThunk } from '../slices/user';
import { useDispatch, useSelector } from '../store';
import { Header } from '../components/header';

const LoginPage: React.FC = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const error = useSelector((state) => state.user.error);

    // Обработчик события нажатия на кнопку "Войти"
    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        const formData = {
            username: (document.getElementById('username') as HTMLInputElement).value,
            password: (document.getElementById('password') as HTMLInputElement).value,
        };

        if (formData.username && formData.password) {
            dispatch(loginThunk(formData)).then((res) => {
                res.meta.requestStatus === 'fulfilled' &&
                    navigate(`${ROUTES.TURBINES_LIST}`); // переход на страницу услуг
            });
        }
    };

    return (
        <div>
            <Header />
            <Container style={{ maxWidth: '400px', marginTop: '150px' }}>
                <h2 style={{ textAlign: 'center', marginBottom: '20px' }}>Вход в VETRYAKI</h2>
                {error && <Alert variant="danger">{error}</Alert>}
                <Form onSubmit={handleSubmit}>
                    <Form.Group controlId="username" style={{ marginBottom: '15px' }}>
                        <Form.Label>Имя пользователя</Form.Label>
                        <Form.Control
                            required
                            type="text"
                            name="username"
                            placeholder="Введите имя пользователя"
                        />
                    </Form.Group>
                    <Form.Group controlId="password" style={{ marginBottom: '20px' }}>
                        <Form.Label>Пароль</Form.Label>
                        <Form.Control
                            required
                            type="password"
                            name="password"
                            placeholder="Введите пароль"
                        />
                    </Form.Group>
                    <Button variant="primary" type="submit" style={{ width: '100%' }}>
                        Войти
                    </Button>
                </Form>
            </Container>
        </div>
    );
};

export default LoginPage;