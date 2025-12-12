import type { FC } from "react"
import { Container, Nav, Navbar } from "react-bootstrap"
import { NavLink, useNavigate } from "react-router"
import { ROUTE_LABELS, ROUTES } from "../routes"
import { logoutThunk } from "../slices/user"
import { useDispatch, useSelector } from "../store"

export const Header: FC<{ mode?: "normal" | "main-page" }> = ({ mode = "normal" }) => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const { isAuthenticated, username } = useSelector(state => state.user);

    // Обработчик события нажатия на кнопку "Выйти"
    const handleExit = async () => {
        await dispatch(logoutThunk());

        navigate(ROUTES.TURBINES_LIST); // переход на страницу списка услуг
    }

    return (
        <header style={{
            borderBottom: mode === "normal" ? "5px solid #5BA1D4" : undefined,
            marginBottom: 24
        }}>
            <Navbar expand="lg">
                <Container>
                    <Navbar.Brand style={{
                        fontSize: 36,
                        lineHeight: "48px",
                        fontWeight: "bold",
                        color: mode === "normal" ? undefined : "white",
                    }}>VETRYAKI</Navbar.Brand>
                    <Navbar.Toggle aria-controls="navbar" />
                    <Navbar.Collapse id="navbar" style={{ justifyContent: "end" }}>
                        <Nav style={{ alignItems: "center", gap: 20 }}>
                            {isAuthenticated
                                ?
                                <div style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>
                                    <NavLink to={ROUTES.PROFILE} style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>{username}</NavLink>
                                    {" "}(<NavLink to={ROUTES.TURBINES_LIST} style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }} onClick={handleExit}>Выйти</NavLink>)
                                </div>
                                : <div style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>
                                    <NavLink to={ROUTES.LOGIN} style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>Войти</NavLink>
                                    {" / "}
                                    <NavLink to={ROUTES.REGISTRATION} style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>Зарегистрироваться</NavLink>
                                </div>
                            }
                            <NavLink to={ROUTES.TURBINES_LIST} style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>{ROUTE_LABELS.TURBINES_LIST}</NavLink>
                            {isAuthenticated && <NavLink to={ROUTES.GENERATION_REQUESTS} style={{ color: mode === "normal" ? "black" : "white", textDecoration: "none" }}>{ROUTE_LABELS.GENERATION_REQUESTS}</NavLink>}
                            <NavLink to={ROUTES.HOME} style={{
                                padding: "12px 30px",
                                color: "white",
                                backgroundColor: "#5BA1D4",
                                borderRadius: 6,
                                textDecoration: "none"
                            }}>Домой</NavLink>
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
        </header >
    )
}