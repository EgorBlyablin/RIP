import { useEffect, useState, type FC, type FormEvent } from "react";
import { Button, Col, Container, Form, InputGroup, Row } from "react-bootstrap";
import { fetchTurbines } from "../api/turbines";
import SearchIcon from "../assets/search.svg?react";
import { Breadcrumbs } from "../components/breadcrumbs";
import { TurbineCard } from "../components/turbine-card";
import { ROUTE_LABELS, ROUTES } from "../routes";
import type { Turbine } from "../api/interfaces";
import { fetchRequestStats } from "../api/generation-requests";

export const TurbinesListPage: FC = () => {
    const [turbinesSearchValue, setTurbinesSearchValue] = useState("");
    const [turbines, setTurbines] = useState<Turbine[]>([]);
    const [turbinesCartCount, setTurbinesCartCount] = useState(0);

    useEffect(() => {
        fetchRequestStats().then(stats => stats && setTurbinesCartCount(stats.TurbinesCount));
    }, []);

    useEffect(() => {
        const fetchTurbinesWrapper = async () => {
            const turbinesData = await fetchTurbines(turbinesSearchValue);
            setTurbines(turbinesData || []);
        };

        fetchTurbinesWrapper();
    }, [turbinesSearchValue]);

    // Handle form submission for search
    const handleSearchSubmit = (e: FormEvent) => {
        e.preventDefault();

        const inputField = document.getElementById("titleFilterValue") as HTMLInputElement;
        setTurbinesSearchValue(inputField.value);
    };

    return (
        <Container className="pb-5">
            <Breadcrumbs
                breadcrumbs={[
                    { label: ROUTE_LABELS.HOME, path: ROUTES.HOME },
                    { label: ROUTE_LABELS.TURBINES_LIST, path: ROUTES.TURBINES_LIST },
                ]}
            />
            <div
                style={{
                    display: "flex",
                    justifyContent: "space-between",
                    paddingBlock: 18,
                }}
            >
                <h1
                    style={{
                        textTransform: "uppercase",
                        fontWeight: "bold",
                        color: "#5B5B5B",
                        alignSelf: "center"
                    }}
                >
                    Ветрогенераторы
                    {turbinesSearchValue && <> (поиск: "{turbinesSearchValue}")</>} {/* Use titleFilter here too */}
                </h1>
                <div style={{ display: "flex", gap: 10, alignItems: "center" }}>
                    <Button
                        variant="secondary"
                        style={{ textWrap: "nowrap", fontSize: "1.5rem", paddingBlock: 0, marginRight: "0.5rem" }}
                    >
                        🖩 <span style={{
                            position: "absolute",
                            background: "white",
                            color: "black",
                            borderRadius: "100rem",
                            border: "1px solid black",
                            padding: "0.05rem 0.5rem",
                            lineHeight: 1,
                            fontSize: "0.7rem",
                        }}>{turbinesCartCount}</span>
                    </Button>
                    <Form onSubmit={handleSearchSubmit}>
                        <InputGroup>
                            <Form.Control
                                name="title-filter"
                                id="titleFilterValue"
                                placeholder="Поиск"
                                defaultValue={turbinesSearchValue}
                                maxLength={20}
                            />
                            <Button type="submit" style={{ display: "flex", alignItems: "center", backgroundColor: "#5BA1D4", border: "none" }}>
                                <SearchIcon style={{ width: "1.2rem", height: "1.2rem" }} />
                            </Button>
                        </InputGroup>
                    </Form>
                </div>
            </div>
            <Row className="row-gap-4">
                {(turbines || []).map((turbine) => (
                    <Col key={turbine.id}>
                        <TurbineCard {...turbine} />
                    </Col>
                ))
                }
            </Row>
        </Container>
    );
};
