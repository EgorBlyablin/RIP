import { Button, Card, Col, Form, Row } from "react-bootstrap";
import { Link } from "react-router";
import type { DsTurbine, DsTurbineGenerationRequest } from "../api/Api";

import defaultImage from "../assets/default.jpg";
import { ROUTES } from "../routes";
import { useDispatch, useSelector } from "../store";
import { addTurbineToGenerationRequest } from "../slices/generation-requests-draft";
import { useState, type FC } from "react";


export const TurbineListCard = (turbine: DsTurbine) => {
    const dispatch = useDispatch()

    const { isAuthenticated } = useSelector(state => state.user)

    return (
        <Card
            style={{
                backgroundColor: "#F5F7F7",
                border: "none",
            }}
        >
            <Card.Img
                style={{ height: 250, objectFit: "cover", objectPosition: "top" }}
                variant="top"
                src={turbine.Image || defaultImage}
                alt={`Фото ${turbine.title}`} />
            <Card.Body className="m-3" style={{ color: "#5B5B5B" }}>
                <Card.Title
                    style={{
                        textTransform: "uppercase",
                        fontWeight: "bold",
                    }}
                    className="mb-4"
                >
                    {turbine.title}
                </Card.Title>
                <div className="mb-4">
                    <div className="d-flex align-items-center mb-2">
                        <span className="me-2 fw-bold">Мощность:</span>
                        <span className="flex-grow-1 border-bottom border-dotted"></span>
                        <span className="ms-2">
                            {((turbine.power || 0) / 1000).toPrecision(2)} кВт
                        </span>
                    </div>
                    <div className="d-flex align-items-center">
                        <span className="me-2 fw-bold">Мачта:</span>
                        <span className="flex-grow-1 border-bottom border-dotted"></span>
                        <span className="ms-2">{turbine.height} м</span>
                    </div>
                </div>
                <div style={{ display: "flex", gap: 10 }}>
                    <Link
                        to={`${ROUTES.TURBINES_LIST}/${turbine.id}`}
                        style={{
                            padding: "12px 15px",
                            color: "white",
                            backgroundColor: "#5BA1D4",
                            textTransform: "uppercase",
                            border: "none",
                            borderRadius: "6px",
                            textDecoration: "none",
                            display: "inline-block",
                        }}
                    >
                        Подробнее
                    </Link>
                    {isAuthenticated && <Button
                        onClick={() => dispatch(addTurbineToGenerationRequest(turbine.id as number))}
                        style={{
                            padding: "12px 15px",
                            color: "white",
                            backgroundColor: "#5BA1D4",
                            textTransform: "uppercase",
                            border: "none",
                            borderRadius: "6px",
                            textDecoration: "none",
                            display: "inline-block",
                        }}
                    >
                        В расчет
                    </Button>}
                </div>
            </Card.Body>
        </Card>
    )
}

export const TurbineGenerationRequestCard: FC<{
    turbine: DsTurbineGenerationRequest;
    editable: boolean;
    onDelete: () => void;
    onSubmit: (avgWind: number | undefined, alpha: number | undefined) => void;
}> = ({ turbine, editable, onDelete, onSubmit }) => {
    const [avgWind, setAvgWind] = useState(turbine.avg_velocity)
    const [alpha, setAlpha] = useState(turbine.alpha)

    const normalize = (v: number | null | undefined) =>
        v === null || v === undefined || Number.isNaN(v) ? null : v;

    const formatGeneration = (value?: number) => {
        if (typeof value !== "number" || isNaN(value)) return "н/д";

        if (value < 1_000) {
            return `${value.toFixed(2)} кВт`;            // < 1 МВт
        } else if (value < 1_000_000) {
            return `${(value / 1_000).toFixed(2)} МВт`;  // < 1000 МВт
        } else {
            return `${(value / 1_000_000).toFixed(2)} ГВт`;
        }
    };

    const isChanged =
        normalize(avgWind) !== normalize(turbine.avg_velocity) ||
        normalize(alpha) !== normalize(turbine.alpha);

    return (
        <Form onSubmit={() => onSubmit(avgWind as any, alpha)}>
            <Row className="align-items-center" style={{ backgroundColor: "#F5F7F7", border: "none", marginBottom: 15, padding: 0 }}>
                <Col md={1} style={{ padding: 0 }}>
                    <Card.Img
                        src={turbine.turbine?.Image || defaultImage}
                        style={{
                            height: 100,
                            width: 100,
                            objectFit: "cover",
                            borderRadius: 4,
                        }}
                    />
                </Col>
                <Col md={2} style={{ paddingLeft: 20, fontWeight: "bold", fontSize: 18 }}>
                    <Link to={`${ROUTES.TURBINES_LIST}/${turbine.turbine?.id}`} style={{ color: "black", textDecoration: "none" }}>
                        {turbine.turbine?.title}
                    </Link>
                </Col>
                <Col md={1} style={{ textAlign: "center" }}>
                    <strong>{turbine.turbine?.power && ((turbine.turbine.power || 0) / 1000).toString()} кВт</strong>
                </Col>
                <Col md={1} style={{ textAlign: "center" }}>
                    <strong>{turbine.turbine?.height} м</strong>
                </Col>
                <Col md={2} style={{ textAlign: "center" }}>
                    {editable
                        ? <Form.Control
                            type="number"
                            step={0.1}
                            min={0.1}
                            max={30}
                            defaultValue={turbine.avgVelocity}
                            value={avgWind as any}
                            onChange={(e) => setAvgWind(Number.parseFloat(e.currentTarget.value))}
                            style={{ width: "100%" }}
                            required
                        />
                        : <strong>{turbine.avg_velocity}</strong>}
                </Col>
                <Col md={2} style={{ textAlign: "center" }}>
                    {editable
                        ? <Form.Control
                            type="number"
                            step={0.01}
                            min={0.01}
                            max={10}
                            defaultValue={turbine.alpha}
                            value={alpha}
                            onChange={(e) => setAlpha(Number.parseFloat(e.currentTarget.value))}
                            style={{ width: "100%" }}
                            required
                        />
                        : <strong>{turbine.alpha}</strong>}
                </Col>
                <Col
                    md={3}
                    style={{
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        position: "relative",
                        textAlign: "center"
                    }}
                >
                    <span style={{ fontSize: 24, fontWeight: 600 }}>
                        {turbine.calculated_generation ? formatGeneration(turbine.calculated_generation) : "н/д"}
                    </span>
                    <div
                        style={{
                            position: "absolute",
                            right: 8,
                            top: "50%",
                            transform: "translateY(-50%)",
                            display: "grid",
                            gridTemplateRows: "1fr 1fr",
                            flexDirection: "column",
                            gap: 12
                        }}
                    >
                        {editable && <>
                            <Button variant="danger" size="sm" onClick={onDelete} style={{ borderRadius: 1000, fontSize: 20, lineHeight: 1, padding: "6px 8px" }}>⨯</Button>
                            {(isChanged && !!avgWind && !!alpha) && <Button size="sm" type="submit" onClick={(e) => {
                                e.preventDefault();
                                onSubmit(avgWind, alpha)
                            }} style={{ borderRadius: 1000, fontSize: 20, lineHeight: 1, padding: "6px 8px", background: "#5BA1D4", borderColor: "#5BA1D4" }}>🖫</Button>}
                        </>
                        }
                    </div>
                </Col>
            </Row>
        </Form>
    )
}