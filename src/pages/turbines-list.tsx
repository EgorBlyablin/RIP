import { useEffect, type FC } from "react";
import { Button, Col, Container, Form, InputGroup, Row } from "react-bootstrap";
import SearchIcon from "../assets/search.svg?react";
import { Breadcrumbs } from "../components/breadcrumbs";
import { TurbineListCard } from "../components/turbine-card";
import { ROUTE_LABELS, ROUTES } from "../routes";
import { getTurbinesList, setTurbinesFilter } from "../slices/turbines";
import { useDispatch, useSelector } from "../store";

import CalculatorIcon from "../assets/calculator.svg?react";
import { useNavigate } from "react-router";
import { BeatLoader } from "react-spinners";

export const TurbinesListPage: FC = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const { isAuthenticated } = useSelector((state) => state.user);
    const { searchFilter, turbines, loading } = useSelector((state) => state.turbines);
    const { generationRequestId, turbinesCount } = useSelector((state) => state.generationRequestDraft);

    useEffect(() => {
        dispatch(getTurbinesList());
    }, []);

    return (
        <Container className="pb-5">
            <Breadcrumbs
                breadcrumbs={[
                    { label: ROUTE_LABELS.HOME, path: ROUTES.HOME },
                    { label: ROUTE_LABELS.TURBINES_LIST, path: ROUTES.TURBINES_LIST },
                ]}
            />
            <Row style={{ display: "flex", gap: 16, justifyContent: "space-between" }} className="p-2 mb-1"            >
                <Col md={"auto"} xs={12}>
                    <h1 style={{ textTransform: "uppercase", fontWeight: "bold", color: "#5B5B5B", lineHeight: 1, margin: 0 }}                    >
                        Ветрогенераторы {searchFilter && <> (поиск: "{searchFilter}")</>}
                    </h1>
                </Col>
                <Col md={"auto"} style={{ display: "flex", gap: 10, alignItems: "center" }}>
                    <Button onClick={() => navigate(`${ROUTES.GENERATION_REQUESTS}/${generationRequestId}`)} variant={!!turbinesCount ? "primary" : "secondary"} disabled={!isAuthenticated || !turbinesCount} style={{ textWrap: "nowrap", fontSize: "1.5rem", height: "36px", marginRight: "0.5rem", display: "flex", alignItems: "center", background: turbinesCount ? "#5BA1D4" : undefined, border: "none" }} >
                        <CalculatorIcon width={"25px"} />
                        <span style={{ position: "absolute", background: "white", color: "black", borderRadius: "100rem", border: "1px solid black", padding: "0.05rem 0.5rem", lineHeight: 1, fontSize: "0.7rem", transform: "translate(100%, -100%)" }}>
                            {turbinesCount || 0}
                        </span>
                    </Button>
                    <Form onSubmit={(e) => {
                        e.preventDefault();
                        dispatch(setTurbinesFilter((document.getElementById("title-filter") as HTMLInputElement | null)?.value))
                        dispatch(getTurbinesList())
                    }}>
                        <InputGroup>
                            <Form.Control id="title-filter" placeholder="Поиск" maxLength={20} key={searchFilter} defaultValue={searchFilter} />
                            <Button type="submit" disabled={loading} style={{ display: "flex", alignItems: "center", backgroundColor: "#5BA1D4", border: "none" }}>
                                <SearchIcon style={{ width: "1.2rem", height: "1.2rem" }} />
                            </Button>
                        </InputGroup>
                    </Form>
                </Col>
            </Row>
            <Row>
                {loading
                    ? <Container style={{display: "flex", justifyContent:"center", paddingTop: "25%"}}>
                        <BeatLoader color="#5BA1D4" />
                    </Container>
                    : (turbines).map((turbine) => (
                        <Col key={turbine.id} className="p-2" xxl={3} lg={4} sm={6} xs={12}>
                            <TurbineListCard {...turbine} />
                        </Col>
                    ))
                }
            </Row>
        </Container>
    );
};
