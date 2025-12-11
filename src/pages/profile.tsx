import { Alert, Button, Card, Container, Form } from "react-bootstrap";
import { useSelector } from "../store";
import { api } from "../api";
import { useState } from "react";
import { Breadcrumbs } from "../components/breadcrumbs";
import { ROUTE_LABELS, ROUTES } from "../routes";

interface AlertData {
    status: "success" | "error" | null;
    message: string
}

export const ProfilePage = () => {
    const user = useSelector(state => state.user);
    const [alert, setAlert] = useState<AlertData>({
        status: null,
        message: ''
    })

    const handlePasswordChange = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const form = e.currentTarget;
        const newPassword = (form.elements.namedItem("newPassword") as HTMLInputElement).value;
        const confirmNewPassword = (form.elements.namedItem("confirmNewPassword") as HTMLInputElement).value;

        if (newPassword !== confirmNewPassword) {
            setAlert({
                status: "error",
                message: "Пароли не совпадают"
            });
            return;
        }

        setAlert({ status: null, message: "" });

        // Здесь можно добавить логику для изменения пароля, например, отправить запрос на сервер
        api.usersUpdate({ password: newPassword }).then(() => {
            setAlert({ status: "success", message: "Пароль успешно изменен" });
            setTimeout(() => setAlert({ status: null, message: "" }), 3000)
        })
    }

    return (
        <Container>
            <Breadcrumbs breadcrumbs={[
                { label: ROUTE_LABELS.HOME, path: ROUTES.HOME },
                { label: ROUTE_LABELS.PROFILE, path: ROUTES.PROFILE },
            ]} />
            <Card className="mt-4 p-5">
                {alert.status && <Alert variant={alert.status === "error" ? "danger" : "success"}>{alert.message}</Alert>}
                <h1>Профиль пользователя {user.username}</h1>
                <hr />
                <Form onSubmit={(e) => handlePasswordChange(e)}>
                    <Form.Group className="mb-3" controlId="newPassword">
                        <Form.Label>Новый пароль</Form.Label>
                        <Form.Control type="password" placeholder="Введите новый пароль"
                            minLength={8} required />
                    </Form.Group>
                    <Form.Group className="mb-3" controlId="confirmNewPassword">
                        <Form.Label>Новый пароль (повторно)</Form.Label>
                        <Form.Control type="password" placeholder="Подтвердите новый пароль"
                            minLength={8} required />
                    </Form.Group>
                    <Button type="submit" style={{ backgroundColor: "#5BA1D4", borderColor: "#5BA1D4" }} >Сохранить изменения</Button>
                </Form>
            </Card>
        </Container>
    )
};